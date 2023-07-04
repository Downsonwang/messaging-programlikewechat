/*
 * @Descripttion:  GroupChat Feature
 * @Author: Downson
 * @Date: 2023-05-08 08:21:37
 * @LastEditTime: 2023-05-14 12:01:18
 */
package models

import "gorm.io/gorm"

//群信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
