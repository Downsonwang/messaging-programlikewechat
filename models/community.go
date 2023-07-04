/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-06-18 10:32:04
 * @LastEditTime: 2023-07-04 08:37:36
 */
package models

import (
	"chatapp/utils"
	"fmt"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}

	if community.OwnerId == 0 {
		return -1, "请先登录"
	}

	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	data := make([]*Community, 10)
	utils.DB.Where("owner_id = ?", ownerId).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data, "查询成功"
}
