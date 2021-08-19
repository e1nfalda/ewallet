package services

import (
	"ewallet/models"

	"github.com/beego/beego/v2/core/logs"
)

const ORDER_STATUS_CREATED = 0
const ORDER_STATUS_FINISHED = 1

/*
 * two step transcation.
 * As far as I konw, most of payment lately is two step. Can extracts use pin code, sms, email etc.
 */

// GetTransactionList get user's transaction list
func GetTransactionList(userPhone string) []models.Transaction {
	return models.GetTransactionList(userPhone)
}

// CreateOrder create a new order for transaction.
func CreateOrder(sender, receiver string, amount float64) (orderId string, errCode int) {
	senderInfo := models.GetUserInfo(sender)
	if senderInfo.Balance < amount {
		errCode = ERROR_CODE_TRANS_1
		return
	}
	orderID, err := models.CreateOrder(sender, receiver, amount)
	if err != nil {
		logs.Error("Error Create Order", sender, receiver, amount, err)
		return "", ERROR_CODE_TRANS_2
	}

	return orderID, ERROR_CODE_SUCCESS
}

func ConfirmOrder(orderID, confirmPin string) int {
	orderInfo := models.GetOrder(orderID)
	if orderInfo == nil {
		return ERROR_CODE_TRANS_3
	}

	if orderInfo.Status != ORDER_STATUS_CREATED {
		return ERROR_CODE_TRANS_6
	}

	userInfo := models.GetUserInfo(orderInfo.FromUser)
	if !userInfo.CheckConfirmPin(confirmPin) {
		return ERROR_CODE_TRANS_4
	}

	if !models.ProcessOrder(orderID) {
		return ERROR_CODE_TRANS_5
	}

	return ERROR_CODE_SUCCESS
}
