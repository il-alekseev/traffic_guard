package keycloakclient

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	u := User{
		Username:   "vasya2",
		FirstName:  "Василий",
		LastName:   "Попов",
		Patronymic: "Петрович",
		Email:      "vas2@local",
		Role:       "SA",
	}
	p := "123456"
	id, err := testClient.CreateUser(u, nil, p)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	t.Log(id)
}

func TestGetUsers(t *testing.T) {
	// users, _ := testClient.GetUsers()
	// for _, user := range users {
	// 	t.Log(user)
	// }
}

func TestGetUserByID(t *testing.T) {
	user, err := testClient.GetUserByID("ec467916-e4f1-4f1e-889b-3ffb92a13388")
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	t.Log(user)

}
