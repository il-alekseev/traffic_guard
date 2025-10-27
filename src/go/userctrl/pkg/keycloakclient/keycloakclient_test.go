package keycloakclient

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

var url string = "https://localhost:8443"
var login string = "admin"
var password string = "admin"
var realm string = "tsum"
var l slog.Logger = *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

var testClient *KeycloakClient = func() *KeycloakClient {
	kc, err := NewKeycloakClient(url, login, password, "master", realm, "grafana-sso", "", "", l)
	slog.Error("NewKeycloakClient", "err", err.Error())
	return kc
}()

func TestToken(t *testing.T) {
	t.Log(testClient.token)
	e, _ := testClient.client.GetCerts(testClient.ctx, testClient.realm)
	t.Log(e)
}

func TestCreateContext(t *testing.T) {
	// 1. Создаем роли для контекста
	//------------------------------
	contextName := "api_ctx"
	caPrefix := "CA"
	coPrefix := "CO"
	// Администратор контекста
	caRole := strings.Join([]string{caPrefix, contextName}, "-")
	// Оператор контекста
	coRole := strings.Join([]string{coPrefix, contextName}, "-")
	caID, err := testClient.CreateRole(caRole)
	if err != nil {
		t.Fatalf("CreateRole(CA): %v", err)
	}
	t.Logf("Success creating role: %s", caID)
	coID, err := testClient.CreateRole(coRole)
	if err != nil {
		t.Fatalf("CreateRole(CO): %v", err)
	}
	t.Logf("Success creating role: %s", coID)
	// 2. Создаем группы для ролей
	//----------------------------
	contextGroupID, err := testClient.CreateGroup(contextName, []string{caPrefix, coPrefix})
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	t.Logf("Success creating group: %s", contextGroupID)
	// 3. Добавляем роли к дочерним группам контекста
	caGroupPath := fmt.Sprintf("/%s/%s", contextName, caPrefix)
	coGroupPath := fmt.Sprintf("/%s/%s", contextName, coPrefix)
	err = testClient.AddRoleToGroupByPath(caGroupPath, caRole)
	if err != nil {
		t.Fatalf("AddRoleToGroupByPath(CA): %v", err)
	}
	err = testClient.AddRoleToGroupByPath(coGroupPath, coRole)
	if err != nil {
		t.Fatalf("AddRoleToGroupByPath(CO): %v", err)
	}
	// 4. Создаем пользователей
	caUser := User{
		Username:   "nagibator999",
		FirstName:  "Василий",
		LastName:   "Попов",
		Patronymic: "Петрович",
		Email:      "nagibator999@test",
		Role:       caRole,
	}
	coUser := User{
		Username:   "antoshka",
		FirstName:  "Антон",
		LastName:   "Головкин",
		Patronymic: "Гаврилович",
		Email:      "antoshka@test",
		Role:       coRole,
	}
	caUserID, err := testClient.CreateUser(caUser, &caGroupPath, "123456")
	if err != nil {
		t.Fatalf("CreateUser(CA): %v", err)
	}
	t.Logf("Success creating user: %s", caUserID)
	coUserID, err := testClient.CreateUser(coUser, &coGroupPath, "123456")
	if err != nil {
		t.Fatalf("CreateUser(CO): %v", err)
	}
	t.Logf("Success creating user: %s", coUserID)
}

func TestCreateSuperAdmin(t *testing.T) {
	// Суперадминистратор
	// 1. Создаем роль
	saRole := "SA"
	saID, err := testClient.CreateRole(saRole)
	if err != nil {
		t.Fatalf("CreateRole(SA): %v", err)
	}
	t.Logf("Success creating role: %s", saID)
	// 2. Создаем группу для роли
	saGroupID, err := testClient.CreateGroup(saRole, nil)
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	t.Logf("Success creating group: %s", saGroupID)
	// 3. Добавляем роли к группе
	err = testClient.AddRoleToGroup(saGroupID, saRole)
	if err != nil {
		t.Fatalf("AddRoleToGroup(SA): %v", err)
	}
	// 4. Создаем пользователей
	saUser := User{
		Username:   "supernova",
		FirstName:  "Феофан",
		LastName:   "Грек",
		Patronymic: "Петрович",
		Email:      "supernova@test",
		Role:       saRole,
	}
	saGroupPath := fmt.Sprintf("/%s", saRole)
	saUserID, err := testClient.CreateUser(saUser, &saGroupPath, "123456")
	if err != nil {
		t.Fatalf("CreateUser(CA): %v", err)
	}
	t.Logf("Success creating user: %s", saUserID)
}

func TestDeleteContext(t *testing.T) {
	// 1. Удаляем роли контекста
	//------------------------------
	contextName := "api_ctx"
	caPrefix := "CA"
	coPrefix := "CO"
	// Администратор контекста
	caRole := strings.Join([]string{caPrefix, contextName}, "-")
	// Оператор контекста
	coRole := strings.Join([]string{coPrefix, contextName}, "-")
	err := testClient.DeleteRoleByName(caRole)
	if err != nil {
		t.Fatalf("DeleteRoleByName(CA): %v", err)
	}
	t.Logf("Success deleting role(CA)")
	err = testClient.DeleteRoleByName(coRole)
	if err != nil {
		t.Fatalf("DeleteRoleByName(CO): %v", err)
	}
	t.Logf("Success deleting role(CO)")
	// 2. Удаляем группы для ролей
	//----------------------------
	path := fmt.Sprintf("/%s", contextName)
	err = testClient.DeleteGroupByPath(path)
	if err != nil {
		t.Fatalf("DeleteGroupByPath: %v", err)
	}
	t.Logf("Success deleting group")
}
