package services

import "encoding/json"

type Result struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// error code defines
const ERROR_CODE_SUCCESS = 0
const ERROR_CODE_USER_1 = 1001 // user not exists
const ERROR_CODE_USER_2 = 1002 // error password
const ERROR_CODE_USER_3 = 1003 // account not exists
const ERROR_CODE_USER_4 = 1004 // need login

const ERROR_CODE_TRANS_1 = 2001 // error amount
const ERROR_CODE_TRANS_2 = 2002 // error create order
const ERROR_CODE_TRANS_3 = 2003 // error get order info error
const ERROR_CODE_TRANS_4 = 2004 // error check confirm pin code error
const ERROR_CODE_TRANS_5 = 2005 // error
const ERROR_CODE_TRANS_6 = 2006 // error order finished

func (p *Result) Json() string {
	data, err := json.Marshal(p)
	if err != nil {
		data, _ = json.Marshal(Result{Status: -1})
		return string(data)
	}
	return string(data)
}
