/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-03 09:25:39
 * @LastEditTime: 2023-07-08 11:03:28
 */
package main

import (
	"chatapp/routers"
	"chatapp/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	utils.InitTimer()
	r := routers.Router()
	r.Run(":8081")
}
