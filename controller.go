package main

import (
	"errors"
	"ewallet/defines"
	"ewallet/services"
	"ewallet/sessions"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

const XSR_SECRET = "oxnSWG75EHmkKrZ7iAZq"

type MainController struct {
	web.Controller
	SessionID string
}

func (p *MainController) Return(err error, data interface{}) {
	var code int
	var message string

	logs.Info("------", err)
	if err == nil {
		code = defines.SUCESS_STATUS
		message = "sucess"
	} else {
		if ewalleErr, ok := err.(*defines.EWALLET_ERROR); ok {
			// user's readable error
			code = ewalleErr.Code
			message = ewalleErr.Desc
		} else {
			// system error
			code = defines.SystemErrorCode
			message = "system error"
			logs.Error("system error:", err)
		}
	}

	rst := &defines.Result{
		Status: code,
		Msg:    message,
		Body:   data,
	}
	p.Ctx.WriteString(rst.Json())
}

func urlNeedValidate(url string) bool {
	noNeeds := map[string]interface{}{
		"/":           0,
		"/user/login": 0,
		"/index":      0,
	}

	_, ok := noNeeds[url]
	return !ok
}

// Prepare overwrite web.Controller method.
func (p *MainController) Prepare() {
	// no need login, just go to next
	if !urlNeedValidate(p.Ctx.Request.URL.Path) {
		return
	}

	sessionID, exists := p.GetSecureCookie(XSR_SECRET, "session_token") // beego implements HMAC encryption.
	logs.Debug("logined info:", sessionID, exists)
	if exists {
		p.SessionID = sessionID
	} else {
		// redirect
		p.Return(defines.ERROR_CODE_USER_1, nil)
	}
}

/******************Controllers******************/

// LoginPage get front html page
func (p *MainController) LoginPage() {
	p.TplName = "index.html"
	p.Render()
}

// MainPage Redirect to login page
func (p *MainController) MainPage() {
	p.Ctx.Redirect(302, "/index")
}

// Login user login
func (p *MainController) Login() {
	phone := p.GetString("phone")
	password := p.GetString("password")

	sessionID, userInfo, err := services.Login(phone, password)
	if err == nil {
		p.SetSecureCookie(XSR_SECRET, "session_token", sessionID)
	}

	p.Return(err, userInfo.FullInfo())
}

// GetAccountBasicInfo user get userself's info
func (p *MainController) GetAccountBasicInfo() {
	phone := p.GetString("phone")
	userInfo, err := services.GetUserInfo(phone)

	p.Return(err, userInfo.PublicInfo())
}

// GetAccountFullInfo get user self's full info
func (p *MainController) GetAccountFullInfo() {
	phone, err := sessions.GetInfo(p.SessionID, "Phone", 1)
	if err != nil {
		p.Return(defines.ERROR_CODE_USER_5, nil)
		return
	}
	userInfo, err := services.GetUserInfo(phone.(string))

	p.Return(err, userInfo.FullInfo())
}

func (p *MainController) CreateTransactionOrder() {
	receiver := p.GetString("toUser")
	amount, err := p.GetFloat("amount")

	if amount <= 0 || err != nil {
		p.Return(defines.ERROR_CODE_TRANS_1, nil)
		return
	}

	sender, err := sessions.GetInfo(p.SessionID, "Phone", 1)
	if err != nil {
		// TODO send alert and redirect to login page
		p.Return(errors.New(""), nil) // error code
		return
	}

	orderID, err := services.CreateOrder(sender.(string), receiver, amount)

	p.Return(err, map[string]string{"order_no": orderID})
}

//  ConfirmTransaction
func (p *MainController) ConfirmTransaction() {
	orderID := p.GetString("orderId")
	confirmPin := p.GetString("pin")

	err := services.ConfirmOrder(orderID, confirmPin)

	p.Return(err, map[string]string{"order_no": orderID})
}

// GetTransactionList 获取历史交易记录
func (p *MainController) GetTransactionList() {
	sender, err := sessions.GetInfo(p.SessionID, "Phone", 1)
	if err != nil {
		p.Return(err, nil)
		return
	}
	transactions := services.GetTransactionList(sender.(string))

	p.Return(err, map[string]interface{}{"transactions": transactions})
}

func main() {
	// web.SetStaticPath("/static", "public")
	web.Router("/", &MainController{}, "get:MainPage")
	web.Router("/index", &MainController{}, "get:LoginPage")
	// User Apis
	web.Router("/user/login", &MainController{}, "post:Login")
	web.Router("/user/get_account_info", &MainController{}, "post:GetAccountBasicInfo")
	web.Router("/user/get_self_info", &MainController{}, "post:GetAccountFullInfo")
	// Transaction Apis
	web.Router("/transaction/create_order", &MainController{}, "post:CreateTransactionOrder")
	web.Router("/transaction/confirm_order", &MainController{}, "post:ConfirmTransaction")
	web.Router("/transaction/list", &MainController{}, "get:GetTransactionList")

	web.Run()
}
