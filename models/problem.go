package models

import (
	"time"

	"gorm.io/gorm"
)

// Problem 作業
type Problem struct {
	gorm.Model
	Class_ID   uint      `gorm:"NOT NULL;"` // 課程 ID , 是 Class.ID
	Problem_ID uint      `gorm:"NOT NULL;"` // 題目 ID
	Start_time time.Time `gorm:"NOT NULL;"` // 作業開始時間
	End_time   time.Time `gorm:"NOT NULL;"` // 作業結束時間
}

//CreateProblem 創建作業
func CreateProblem(problem *Problem) (err error) {
	err = DB.Create(&problem).Error
	return
}

//UpdateProblem 更新作業
func UpdateProblem(problem *Problem) (err error) {
	err = DB.Where("id = ?", problem.ID).Save(&problem).Error
	return
}

//DeleteProblem 刪除作業
func DeleteProblem(id uint) (err error) {
	err = DB.Delete(&Problem{}, id).Error
	return
}

//ListProblem 列出所有作業
func ListProblem(class_id uint) (problem []Problem, err error) {
	err = DB.Where("class_id =?", class_id).Find(&problem).Error
	return
}

//ProblemByProblemID 用 作業 ID 查詢作業
func ProblemByProblemID(id uint) (Problem, error) {
	var problem Problem

	if err := DB.Where("id = ?", id).First(&problem).Error; err != nil {

		return Problem{}, err
	}

	return problem, nil
}
