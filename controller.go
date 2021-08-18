package main

import (
	"ewallet/services"

	"github.com/beego/beego/v2/server/web"
)

const XSR_SECRET = "oxnSWG75EHmkKrZ7iAZq"

// error code defines
const ERROR_CODE_SUCCESS = 0
const ERROR_CODE_USER_1 = 1001 // user not exists
const ERROR_CODE_USER_2 = 1002 // error password
const ERROR_CODE_USER_3 = 1003 // account not exists
const ERROR_CODE_USER_4 = 1004 // need login

type MainController struct {
	web.Controller
	SessionID string
}

func (p *MainController) Return(status int, data interface{}) {
	rst := &Result{
		status: status,
		body:   data,
	}
	p.Ctx.WriteString(rst.Json())
}

// Prepare overwrite web.Controller method.
func (p *MainController) Prepare() {
	if p.Ctx.Request.URL.Path == "/usr/login" {
		return
	}
	sessionID, exists := p.GetSecureCookie(XSR_SECRET, "session_token") // beego implements HMAC encryption.
	if exists {
		p.SessionID = sessionID
	} else {
		p.Return(ERROR_CODE_USER_1, nil)
		p.Ctx.WriteString("go to regiter page...")
	}
}

func (p *MainController) Login() {
	phone := p.GetString("phone")
	password := p.GetString("password")

	sessionID, err := services.Login(phone, password)
	if err != nil {
		p.Return(ERROR_CODE_USER_1, nil)
	}
	userInfo := services.GetUserInfo()
	p.SetSecureCookie(XSR_SECRET, "session_token", sessionID)

	p.Return(ERROR_CODE_SUCCESS, userInfo)
}

func (p *MainController) LoginPage() {
}

func main() {
	web.Router("/", &MainController{}, "get:LoginPage")
	// User Apis
	web.Router("/user/login", &MainController{}, "get:Login")
	web.Router("/user/get_account_intro", &MainController{}, "get:GetAccountIntro")

	// Transaction Apis
	web.Router("/transaction/list", &MainController{}, "get:GetTransactionList")
	web.Router("/transaction/detail", &MainController{}, "get:GetTransactionDetail")
	web.Router("/transaction/create_order", &MainController{}, "post:CreateTransactionOrder")
	web.Router("/transaction/confirm", &MainController{}, "post:ConfirmTransaction")

	web.Run()
}

func init() {
}
