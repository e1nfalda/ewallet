package services

import (
	"crypto/sha256"
	"ewallet/defines"
	"ewallet/models"
	"ewallet/sessions"
	"fmt"
)

// GetUserInfo  get user's information。
func GetUserInfo(phone string) (userInfo *models.User, err error) {
	userInfo = models.GetUserInfo(phone)
	if userInfo == nil {
		err = defines.ERROR_CODE_USER_1
		return
	}
	return
}

// Login check user phone and password.If success return sessionID
func Login(phone, password string) (sessionID string, userInfo *models.User, err error) {
	userInfo = models.GetUserInfo(phone)
	if userInfo == nil {
		err = defines.ERROR_CODE_USER_1
		return
	}
	if !checkPassword(userInfo, password) {
		userInfo = nil
		err = defines.ERROR_CODE_USER_2
		return
	}
	sessionID = sessions.CreateSession(phone)
	return
}

func checkPassword(userInfo *models.User, password string) bool {
	// logs.Debug("check password. encryption:", fmt.Sprintf("%x", sha256.Sum256([]byte(password+userInfo.Salt))), userInfo.Password)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+userInfo.Salt))) == userInfo.Password
}
