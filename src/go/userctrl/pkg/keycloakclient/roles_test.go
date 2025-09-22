package keycloakclient

import "testing"

func TestCreateRole(t *testing.T) {
	id, err := testClient.CreateRole("ct2-CA")
	if err != nil {
		t.Fatalf("CreateRole: %v", err)
	}
	t.Log(id)
}
