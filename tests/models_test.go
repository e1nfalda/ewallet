package tests

import (
	"ewallet/models"
	"testing"
)

// TestUserInfoOrm test database service and demo data inputed
// database should be deployed.
func TestUserInfoOrm(t *testing.T) {
	userInfo := models.GetUserInfo("123")
	if userInfo == nil {
		t.Error("Error Get UserInfo")
	}

	// models.CreateOrder("123", "234")

}
