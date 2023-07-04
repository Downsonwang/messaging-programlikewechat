/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-04 08:27:01
 * @LastEditTime: 2023-06-04 16:13:12
 */
package service

import (
	"chatapp/models"
	"fmt"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex GetIndex
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags HomeWebPage
// @Accept application/json
// @Produce application/json
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	ind.Execute(c.Writer, "index")
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

// 跳转到聊天页面
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/main.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
	)
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	ind.Execute(c.Writer, user)
}
