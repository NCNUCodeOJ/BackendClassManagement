package models

import (
	"gorm.io/gorm"
)

// Class 課程
type Class struct {
	gorm.Model
	Class_Name string `gorm:"type:text;"` // 課程名稱
	Teacher    uint   `gorm:"NOT NULL;"`  // 老師

}

//CreateClass 創建課程
func CreateClass(class *Class) (err error) {
	err = DB.Create(&class).Error
	return
}

//UpdateClass 更新課程
func UpdateClass(class *Class) (err error) {
	err = DB.Where("id =?", class.ID).Save(&class).Error
	return
}

//DeleteClass 刪除課程
func DeleteClass(id uint) (err error) {
	err = DB.Delete(&Class{}, id).Error
	return
}

//ListClass 列出所有課程
func ListClass() (class []Class, err error) {
	err = DB.Find(&class).Error
	return
}

//ClassByClassID 用 課程 ID 查詢課程
func ClassByClassID(id uint) (Class, error) {
	var class Class

	if err := DB.Where("id = ?", id).First(&class).Error; err != nil {

		return Class{}, err
	}

	return class, nil
}
