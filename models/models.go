package models

/*
 * models
 In beego v2.x, it’s thread safe.

*/

import (
	"errors"
	"ewallet/defines"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

var DB_HOST = os.Getenv("EWALLET_DBHOST")

const DATABASE = "ewallet"
const DB_USER = "root"
const DB_PASSWORD = "123456"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type User struct {
	Id          int
	PhoneNo     string    `orm:"size(20);unique"`
	Name        string    `orm:"size(20)"`
	Password    string    `orm:"size(64)"`
	Salt        string    `orm:"size(10)"`
	Balance     float64   `orm:"digits(12);decimals(2)"`
	TransterPin string    `orm:"size(64)"`
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateDate  time.Time `orm:"auto_now;type(datetime)"`
}

type Transaction struct {
	Id         int
	OrderID    string  `orm:"size(15);unique"`
	FromUser   string  `orm:"size(20)"`
	ToUser     string  `orm:"size(20)"`
	Amount     float64 `orm:"digits(12);decimals(2)"`
	Status     int
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateDate time.Time `orm:"auto_now;type(datetime)"`
}

func (u *User) TableName() string {
	return "auth_user"
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		{"PhoneNo"},
	}
}

// FullInfo Maily for user self user or higher level user use
func (u *User) FullInfo() map[string]interface{} {
	return map[string]interface{}{
		"Name":       u.Name,
		"PhoneNo":    u.PhoneNo,
		"Balance":    u.Balance,
		"updateDate": u.UpdateDate,
	}
}

// PublicInfo Other user can see info.To Avoid sesstive info
func (u *User) PublicInfo() map[string]interface{} {
	return map[string]interface{}{
		"Name":    fmt.Sprintf("%c*%c", u.Name[0], u.Name[len(u.Name)-1]),
		"PhoneNo": u.PhoneNo,
	}
}

// CheckConfirmPin check confirm pin is correct
func (u *User) CheckConfirmPin(inputPin string) bool {
	return u.TransterPin == inputPin
}

func (u *Transaction) TableName() string {
	return "transaction"
}

func (u *Transaction) TableIndex() [][]string {
	return [][]string{
		{"FromUser"},
		{"ToUser"},
		{"CreateDate"},
	}
}

func CreateOrder(sender, receiver string, amount float64) (string, error) {
	o := orm.NewOrm()
	var order Transaction
	order.OrderID = genOrderID(15) // strict way needs check order exists. to simply emit.
	order.FromUser = sender
	order.ToUser = receiver
	order.Amount = amount

	_, err := o.Insert(&order)
	if err != nil {
		return "", err
	}
	return order.OrderID, nil
}

func GetOrder(orderID string) *Transaction {
	o := orm.NewOrm()
	orderInfo := Transaction{OrderID: orderID}

	err := o.Read(&orderInfo, "orderID")

	if err == orm.ErrNoRows {
		fmt.Println("No result found.")
		return nil
	}

	return &orderInfo
}

// GetTransactionList get phone finished transcations
func GetTransactionList(phone string) []Transaction {
	var trans = make([]Transaction, 0, 10)

	o := orm.NewOrm()

	qs := o.QueryTable("transaction")
	qs = qs.Filter("status", 1)
	cond := qs.GetCond()
	cond = cond.And("from_user", phone).Or("to_user", phone)
	cond = cond.And("status", 1)
	qs = qs.SetCond(cond)
	qs.OrderBy("-id").Limit(10).All(&trans)

	return trans
}

func GetUserInfo(phone string) *User {
	o := orm.NewOrm()
	user := User{PhoneNo: phone}

	err := o.Read(&user, "PhoneNo")

	if err == orm.ErrNoRows {
		fmt.Println("No result found.")
		return nil
	}

	return &user
}

// ProcessOrder Finish a order in database's transcation
func ProcessOrder(orderID string) error {
	o := orm.NewOrm()
	orderInfo := Transaction{OrderID: orderID}

	// start database transcation
	to, err := o.Begin()
	if err != nil {
		logs.Error("Error process order:", err)
		return errors.New(fmt.Sprintf("error start transcation %s", err))
	}

	err = to.ReadForUpdate(&orderInfo, "orderID")
	if err != nil {
		logs.Error("Error process order:", err)
		return errors.New(fmt.Sprintf("error read order info %s", err))
	}
	if orderInfo.Status != defines.ORDER_STATUS_CREATED {
		to.Rollback()
		return defines.ERROR_CODE_TRANS_7
	}

	sender := User{PhoneNo: orderInfo.FromUser}
	err = to.ReadForUpdate(&sender, "PhoneNo")
	if err != nil {
		logs.Error("Error Get sender:", err)
		return errors.New(fmt.Sprintf("error read sender info %s", err))
	}
	if sender.Balance < orderInfo.Amount {
		to.Rollback()
		return defines.ERROR_CODE_TRANS_8
	}

	_, err = to.Raw("UPDATE auth_user SET balance=balance-? where phone_no=?", orderInfo.Amount, orderInfo.FromUser).Exec()
	if err != nil {
		logs.Error("Error process order:", err)
		err = to.Rollback()
		return errors.New(fmt.Sprintf("error deduct sender balance %s", err))
	}
	_, err = to.Raw("UPDATE auth_user SET balance=balance+? where phone_no=?", orderInfo.Amount, orderInfo.ToUser).Exec()
	if err != nil {
		logs.Error("Error erorder:", err)
		to.Rollback()
		return errors.New(fmt.Sprintf("error add receiver balance %s", err))
	}

	orderInfo.Status = defines.ORDER_STATUS_FINISHED
	to.Update(&orderInfo)

	err = to.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("error commit order info %s", err))
	}

	return nil
}

// CreateTables create tables when deploy
func CreateTables() {
	// Database alias.
	name := "default"

	// Drop table and re-create.
	force := false

	// Print log.
	verbose := true

	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8", DB_USER, DB_PASSWORD, DB_HOST, DATABASE))

	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Transaction))

}

// genSessionID simple random algorithm
func genOrderID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterBytes))]
	}
	return string(b)
}
