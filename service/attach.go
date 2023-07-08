/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-06-03 10:33:30
 * @LastEditTime: 2023-07-08 10:08:31
 */
package service

import (
	"chatapp/utils"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	UploadLocal(c)
}

//上传文件到本地
func UploadLocal(c *gin.Context) {
	w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file")
	defer srcFile.Close()
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	suffix := ".png"
	ofileName := head.Filename
	tem := strings.Split(ofileName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "发送图片成功")
}
