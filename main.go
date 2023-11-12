package main

import (
	"dbCLI/cmd"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:MyNewPass@tcp(127.0.0.1:3306)/Lab1?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	cmd.SetDB(db)

	cmd.Execute()

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}
