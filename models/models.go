package models

/*
 * models
 In beego v2.x, it’s thread safe.

*/

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

const DATABASE = "ewallet"
const DB_USER = "root"
const DB_PASSWORD = "123456"

type User struct {
	Id          int
	PhoneNo     string    `orm:"size(20)"`
	Name        string    `orm:"size(20)"`
	Password    string    `orm:"size(20)"`
	Salt        string    `orm:"size(10)"`
	Balance     float64   `orm:"digits(12);decimals(2)"`
	TransterPin string    `orm:"size(10)"`
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateDate  time.Time `orm:"auto_now;type(datetime)"`
}

type Transaction struct {
	Id         int
	FromUser   int
	ToUser     int
	Amount     float32 `orm:"digits(12);decimals(2)"`
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

func GetUserInfo(phone string) {
	o := orm.NewOrm()
	user := User{PhoneNo: phone}

	err := o.Read(&user)

	if err == orm.ErrNoRows {
		fmt.Println("No result found.")
	} else if err == orm.ErrMissPK {
		fmt.Println("No primary key found.")
	} else {
		fmt.Println(user.Id, user.Name)
	}
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", DB_USER, DB_PASSWORD, DATABASE))

	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Transaction))

	createTables()
}

func createTables() {
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
