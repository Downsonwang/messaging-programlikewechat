/*
 * @Descripttion:
 * @Author: IM
 * @Date: 2023-05-08 08:18:03
 * @LastEditTime: 2023-06-04 15:03:49
 */
package models

import (
	"chatapp/utils"
	"fmt"

	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //对应的类型 1是好友 2是群组 3是xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contact := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contact)
	for _, v := range contact {
		objIds = append(objIds, uint64(v.TargetId))
		fmt.Println(v)
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友
func AddFriend(userId uint, targetId uint) (int, string) {
	user := UserBasic{}
	if targetId != 0 {
		user = FindByID(targetId)
		fmt.Println(targetId, "       ", userId)
		if user.Salt != "" {
			if userId == user.ID {
				return -1, "不能添加自己"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetId).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不能重复添加"
			}
			tx := utils.DB.Begin()
			// 发生任何异常都会滚
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contactFir := Contact{}
			contactFir.OwnerId = targetId
			contactFir.TargetId = userId
			contact.Type = 1
			if err := utils.DB.Create(&contactFir).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成长"
		}
		return -1, "没有找到此用户"
	}
	return -1, "好友ID不等于空"
}
