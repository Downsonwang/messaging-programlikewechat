/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-04 08:16:11
 * @LastEditTime: 2023-07-04 21:13:13
 */
package routers

import (
	"chatapp/docs"
	"chatapp/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//加载静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	// r.StaticFS()
	r.LoadHTMLGlob("views/**/*")

	//Home
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/searchFriends", service.SearchFriends)
	//用户模块
	r.GET("/user/getlist", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/attach/upload", service.Upload)
	r.POST("/contact/createCommunity", service.CreateCommunity)
	r.POST("/contact/loadCommunity", service.LoadCommunity)

	r.POST("/contact/joinGroup", service.JoinGroups)
	//创建群聊
	r.POST("", service.AddFriend)
	return r
}
