/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-03 09:25:39
 * @LastEditTime: 2023-05-07 20:17:43
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
	r := routers.Router()
	r.Run(":8081")
}
