package dbstructurer

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type TableInfo struct {
	TableName string
	Columns   []ColumnInfo
}

type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
	IsPrimary  bool
}

func GetDBStructure() {
	// Параметры подключения
	connStr := "host=192.168.130.112 port=5432 user=ksu_user password=ksu_password dbname=ksu sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL")

	// Получаем список таблиц
	tables, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем структуру для каждой таблицы
	for i, table := range tables {
		columns, err := getTableStructure(db, table)
		if err != nil {
			log.Printf("Ошибка получения структуры таблицы %s: %v", table, err)
			continue
		}

		tableInfo := TableInfo{
			TableName: table,
			Columns:   columns,
		}

		printTableInfo(tableInfo)

		// Разделитель между таблицами
		if i < len(tables)-1 {
			fmt.Println("\n" + strings.Repeat("-", 50))
		}
	}
}

func getTables(db *sql.DB) ([]string, error) {
	query := `
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        ORDER BY table_name
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func getTableStructure(db *sql.DB, tableName string) ([]ColumnInfo, error) {
	query := `
        SELECT 
            c.column_name,
            c.data_type,
            c.is_nullable,
            CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary
        FROM information_schema.columns c
        LEFT JOIN (
            SELECT ku.column_name
            FROM information_schema.table_constraints tc
            JOIN information_schema.key_column_usage ku 
                ON tc.constraint_name = ku.constraint_name
            WHERE tc.constraint_type = 'PRIMARY KEY' 
            AND ku.table_name = $1
        ) pk ON c.column_name = pk.column_name
        WHERE c.table_name = $1
        ORDER BY c.ordinal_position
    `

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(&col.ColumnName, &col.DataType, &col.IsNullable, &col.IsPrimary)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func printTableInfo(table TableInfo) {
	fmt.Printf("\n📊 Таблица: %s\n", table.TableName)
	fmt.Printf("%-25s %-20s %-10s %s\n", "Столбец", "Тип", "Null", "PK")
	fmt.Println(strings.Repeat("-", 65))

	for _, col := range table.Columns {
		pk := ""
		if col.IsPrimary {
			pk = "🔑"
		}
		fmt.Printf("%-25s %-20s %-10s %s\n",
			col.ColumnName, col.DataType, col.IsNullable, pk)
	}
}
