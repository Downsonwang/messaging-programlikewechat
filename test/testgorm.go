/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-03 09:34:16
 * @LastEditTime: 2023-06-18 11:07:26
 */
package main

import (
	"chatapp/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Community{})
}
