package models

import "testing"

// TestUserInfoOrm test database service and demo data inputed
// database should be deployed.
func TestUserInfoOrm(t *testing.T) {
	userInfo := GetUserInfo("123")
	if userInfo == nil {
		t.Error("Error Get UserInfo")
	}
}
