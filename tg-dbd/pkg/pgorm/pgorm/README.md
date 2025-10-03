## Пример использования
```go
package main

import (
	"fmt"
	"yourpackage/postgres"
	"time"
	
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"uniqueIndex"`
}

func main() {
	// Строка подключения
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable"
	
	// Создаем подключение с опциями
	pg, err := postgres.New(dsn,
		postgres.MaxPoolSize(20),
		postgres.AutoMigrate(true),
		postgres.Models(&models.Context{}, &models.DefaultDashboard{}),
	)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	// Использование GORM
	var user User
	if err := pg.DB.First(&user).Error; err != nil {
		fmt.Println("No users found")
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	// Создание нового пользователя
	newUser := User{Name: "John Doe", Email: "john@example.com"}
	if err := pg.DB.Create(&newUser).Error; err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Println("User created successfully")
	}

	// Использование транзакции
	err = pg.WithTx(context.Background(), func(tx *gorm.DB) error {
		if err := tx.Create(&User{Name: "Alice", Email: "alice@example.com"}).Error; err != nil {
			return err
		}
		return tx.Create(&User{Name: "Bob", Email: "bob@example.com"}).Error
	})
	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	}
}
```