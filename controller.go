package main

import (
	"ewallet/services"
	"ewallet/sessions"

	// "github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

const XSR_SECRET = "oxnSWG75EHmkKrZ7iAZq"

type MainController struct {
	web.Controller
	SessionID string
}

func (p *MainController) Return(status int, data interface{}) {
	rst := &services.Result{
		Status: status,
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
	if !urlNeedValidate(p.Ctx.Request.URL.Path) {
		return
	}
	sessionID, exists := p.GetSecureCookie(XSR_SECRET, "session_token") // beego implements HMAC encryption.
	logs.Debug("logined info:", sessionID, exists)
	if exists {
		p.SessionID = sessionID
	} else {
		p.Return(services.ERROR_CODE_USER_1, nil)
	}
}

func (p *MainController) Login() {
	phone := p.GetString("phone")
	password := p.GetString("password")

	sessionID, userInfo, errCode := services.Login(phone, password)
	if errCode == services.ERROR_CODE_SUCCESS {
		p.SetSecureCookie(XSR_SECRET, "session_token", sessionID)
	}

	p.Return(errCode, userInfo)
}

func (p *MainController) LoginPage() {
	p.TplName = "index.html"
	p.Render()
}

func (p *MainController) MainPage() {
	p.Ctx.Redirect(302, "/index")
}

func (p *MainController) GetAccountIntro() {
	phone := p.GetString("phone")
	userInfo, errCode := services.GetUserPublicInfo(phone)

	p.Return(errCode, userInfo)
}

func (p *MainController) CreateTransactionOrder() {
	receiver := p.GetString("toUser")
	amount, err := p.GetFloat("amount")

	if amount <= 0 || err != nil {
		p.Return(services.ERROR_CODE_TRANS_1, nil)
		return
	}

	sender, err := sessions.GetInfo(p.SessionID, "Phone", 1)
	logs.Debug("------------", p.SessionID, sender, err)
	if err != nil {
		// TODO send alert and redirect to login page
		p.Return(10000, nil) // error code
		return
	}

	orderID, errCode := services.CreateOrder(sender.(string), receiver, amount)

	p.Return(errCode, map[string]string{"order_no": orderID})
}

func (p *MainController) ConfirmTransaction() {
	orderID := p.GetString("orderId")
	confirmPin := p.GetString("pin")

	errCode := services.ConfirmOrder(orderID, confirmPin)

	p.Return(errCode, map[string]string{"order_no": orderID})
}

func (p *MainController) GetTransactionList() {

	sender, err := sessions.GetInfo(p.SessionID, "Phone", 1)
	if err != nil {
		// TODO send alert and redirect to login page
		p.Return(10000, nil) // error code
		return
	}
	transactions := services.GetTransactionList(sender.(string))

	p.Return(services.ERROR_CODE_SUCCESS, map[string]interface{}{"transactions": transactions})
}

func main() {
	// web.SetStaticPath("/static", "public")
	web.Router("/", &MainController{}, "get:MainPage")
	web.Router("/index", &MainController{}, "get:LoginPage")
	// User Apis
	web.Router("/user/login", &MainController{}, "post:Login")
	web.Router("/user/get_account_intro", &MainController{}, "post:GetAccountIntro")
	// Transaction Apis
	web.Router("/transaction/create_order", &MainController{}, "post:CreateTransactionOrder")
	web.Router("/transaction/confirm_order", &MainController{}, "post:ConfirmTransaction")
	web.Router("/transaction/list", &MainController{}, "get:GetTransactionList")

	web.Run()
}
