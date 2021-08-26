package defines

import (
	"encoding/json"
	"fmt"
)

// order status
const (
	ORDER_STATUS_CREATED  = 0
	ORDER_STATUS_FINISHED = 1
)

// get user info type
const (
	USER_FULL_INFO   = 0
	USER_SIMPLE_INFO = 1
)

const SUCESS_STATUS = 0

type Result struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Body   interface{} `json:"body"`
}

type EWALLET_ERROR struct {
	Code int
	Desc string
}

func (p *EWALLET_ERROR) Error() string {
	return fmt.Sprintf("Error Desc: %s", p.Desc)
}

func (p *EWALLET_ERROR) ErrorCode() int {
	return p.Code
}

func newEwalletError(code int, desc string) *EWALLET_ERROR {
	return &EWALLET_ERROR{code, desc}
}

var (
	SystemErrorCode   = -1
	ERROR_CODE_USER_1 = newEwalletError(1001, "user not exists")
	ERROR_CODE_USER_2 = newEwalletError(1002, "error password")
	ERROR_CODE_USER_3 = newEwalletError(1003, "account not exists")
	ERROR_CODE_USER_4 = newEwalletError(1004, "need login")
	ERROR_CODE_USER_5 = newEwalletError(1004, "need relogin")

	ERROR_CODE_TRANS_1 = newEwalletError(2001, "balance not enough")
	ERROR_CODE_TRANS_2 = newEwalletError(2002, "can't create order")
	ERROR_CODE_TRANS_3 = newEwalletError(2003, "can't get order info")
	ERROR_CODE_TRANS_4 = newEwalletError(2004, "pin code incorrect")
	ERROR_CODE_TRANS_5 = newEwalletError(2005, "transfer error")
	ERROR_CODE_TRANS_6 = newEwalletError(2006, "error order finished")
	ERROR_CODE_TRANS_7 = newEwalletError(2007, "error order status")
	// ERROR_CODE_TRANS_8  = newEwalletError(2008, "sender balance not enough")
	ERROR_CODE_TRANS_9  = newEwalletError(2009, "can not transfer money to yourself")
	ERROR_CODE_TRANS_10 = newEwalletError(20010, "amount cann't be negative")
)

func (p *Result) Json() string {
	data, err := json.Marshal(p)
	if err != nil {
		data, _ = json.Marshal(Result{Status: -1})
		return string(data)
	}
	return string(data)
}
