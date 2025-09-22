package keycloakclient

import "testing"

func TestCreateGroup(t *testing.T) {
	//testClient.CreateGroup("8")
}

func TestGetGroups(t *testing.T) {
	testClient.GetGroups()
}

func TestUpdateGroup(t *testing.T) {
	testClient.UpdateGroup("27104562-cd89-49b8-911c-cfe178252b0c", "New")
}

func TestDeleteGroup(t *testing.T) {
	testClient.DeleteGroupByID("27104562-cd89-49b8-911c-cfe178252b0c")
}
