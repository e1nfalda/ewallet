package main

import (
	"ewallet/models"
	_ "ewallet/models"
	"fmt"
)

func main() {
	fmt.Println("start create tables")
	models.CreateTables()
	fmt.Println("tables created")
}
