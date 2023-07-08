/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-04 08:16:11
 * @LastEditTime: 2023-07-08 10:42:55
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
	r.POST("/user/find", service.FindByID)

	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)
	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	//发送消息
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//群列表
	r.POST("/contact/loadCommunity", service.LoadCommunity)
	//加入群
	r.POST("/contact/joinGroup", service.JoinGroups)

	return r
}
