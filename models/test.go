package models

import (
	"gorm.io/gorm"
)

// Class 課程
type Test struct {
	gorm.Model
	Class_ID     uint   `gorm:"type:text;"` // 課程 ID，是 models.Class.ID
	TestPaper_ID uint   `gorm:"type:text;"` // 測驗卷 ID
	Start_time   string `gorm:"NOT NULL;"`  // 測驗開始時間
	End_time     string `gorm:"NOT NULL;"`  // 測驗結束時間
}

//CreateTest 創建測驗
func CreateTest(test *Test) (err error) {
	err = DB.Create(&test).Error
	return
}

//UpdateTest 更新測驗
func UpdateTest(test *Test) (err error) {
	err = DB.Where("id = ?", test.ID).Save(&test).Error
	return
}

//DeleteTest 刪除測驗
func DeleteTest(id uint) (err error) {
	err = DB.Delete(&Test{}, id).Error
	return
}

//ListTest 列出所有測驗
func ListTest() (test []Test, err error) {
	err = DB.Find(&test).Error
	return
}

//TestkByTestID 用 測驗 ID 查詢測驗
func TestkByTestID(id uint) (Test, error) {
	var test Test

	if err := DB.Where("id = ?", id).First(&test).Error; err != nil {

		return Test{}, err
	}

	return test, nil
}
