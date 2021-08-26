package services

import (
	"ewallet/defines"
	"ewallet/models"
)

/*
 * two step transcation.
 * As far as I konw, most of payment lately is two step. Can extracts use pin code, sms, email etc.
 */

// GetTransactionList get user's transaction list
func GetTransactionList(userPhone string) []models.Transaction {
	return models.GetTransactionList(userPhone)
}

// CreateOrder create a new order for transaction.
func CreateOrder(sender, receiver string, amount float64) (orderId string, err error) {
	if amount < 0 {
		err = defines.ERROR_CODE_TRANS_10
		return
	}
	if sender == receiver {
		err = defines.ERROR_CODE_TRANS_9
		return
	}
	senderInfo := models.GetUserInfo(sender)
	if senderInfo.Balance < amount {
		err = defines.ERROR_CODE_TRANS_1
		return
	}
	return models.CreateOrder(sender, receiver, amount)
}

func ConfirmOrder(orderID, confirmPin string) error {
	orderInfo := models.GetOrder(orderID)
	if orderInfo == nil {
		return defines.ERROR_CODE_TRANS_3
	}
	userInfo := models.GetUserInfo(orderInfo.FromUser)
	if !userInfo.CheckConfirmPin(confirmPin) {
		return defines.ERROR_CODE_TRANS_4
	}

	return models.ProcessOrder(orderID)
}
