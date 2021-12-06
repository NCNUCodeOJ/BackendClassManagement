package models

import (
	"gorm.io/gorm"
)

// ClassUser 課程使用者
type ClassUser struct {
	gorm.Model
	Class_ID uint `gorm:"NOT NULL;"`          // 課程 ID
	User_ID  uint `gorm:"NOT NULL;"`          // 使用者 ID
	Role     int  `gorm:"NOT NULL;default:0"` // 角色

}

//CreateClass 新增課程使用者
func CreateClassUser(classuser *ClassUser) (err error) {
	err = DB.Create(&classuser).Error
	return
}

//UpdateClass 變更課程使用者
func UpdateClassUser(classuser *ClassUser) (err error) {
	err = DB.Where("class_id = ? AND user_id=?", classuser.Class_ID, classuser.User_ID).Save(&classuser).Error
	return
}

//DeleteClass 刪除課程使用者
func DeleteClassUser(user_id uint, class_id uint) (err error) {
	err = DB.Where("user_id =? AND class_id =?", user_id, class_id).Delete(&ClassUser{}).Error
	return
}

//ListClass 列出所有課程使用者
func ListClassUser(class_id uint) (classuser []ClassUser, err error) {
	err = DB.Where("class_id = ?", class_id).Find(&classuser).Error
	return
}

//ListClass 列出課程使用者的所有課堂
func ListClassUserClass(user_id uint) (classuser []ClassUser, err error) {
	err = DB.Where("user_id = ?", user_id).Find(&classuser).Error
	return
}

//ClassUserByClassUserID 用 課程使用者 ID 與 課程 ID 查詢課程使用者
func ClassUserByClassUserID(user_id uint, class_id uint) (ClassUser, error) {
	var classuser ClassUser

	if err := DB.Where("user_id = ? AND class_id = ?", user_id, class_id).First(&classuser).Error; err != nil {

		return ClassUser{}, err
	}

	return classuser, nil
}
