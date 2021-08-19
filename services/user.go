package services

import (
	"crypto/sha256"
	"ewallet/models"
	"ewallet/sessions"
	"fmt"
)

// GetUserPublicInfo
func GetUserPublicInfo(phone string) (userInfo *models.User, errCode int) {
	userInfo = models.GetUserInfo(phone)
	if userInfo == nil {
		errCode = ERROR_CODE_USER_1
		return
	}
	return
}

// Login check user phone and password.If success return sessionID
func Login(phone, password string) (sessionID string, userInfo *models.User, errCode int) {
	userInfo = models.GetUserInfo(phone)
	if userInfo == nil {
		errCode = ERROR_CODE_USER_1
		return
	}
	if !checkPassword(userInfo, password) {
		userInfo = nil
		errCode = ERROR_CODE_USER_2
		return
	}
	sessionID = sessions.CreateSession(phone)
	return
}

func checkPassword(userInfo *models.User, password string) bool {
	// logs.Debug("check password. encryption:", fmt.Sprintf("%x", sha256.Sum256([]byte(password+userInfo.Salt))), userInfo.Password)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+userInfo.Salt))) == userInfo.Password
}
