package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Problem 作業
type Problem struct {
	gorm.Model
	Class_ID   uint      `gorm:"NOT NULL;"`  // 課程 ID , 是 Class.ID
	Problem_ID uint      `gorm:"NOT NULL;"`  // 題目 ID
	Language   string    `gorm:"NOT NULL;"`  // 程式語言
	Start_time time.Time `gorm:"NOT NULL;"`  // 作業開始時間
	End_time   time.Time `gorm:"NOT NULL;"`  // 作業結束時間
	Moss       string    `gorm:"default:"";` // Moss 資料庫
}

// Submission 作業提交
type Submission struct {
	gorm.Model
	ProblemID           uint `gorm:"NOT NULL;"` // 作業 ID
	UserID              uint `gorm:"NOT NULL;"` // 使用者 ID
	PrivateSubmissionID uint `gorm:"NOT NULL;"` // 私有提交
}

//CreateSubmission 創建提交
func CreateSubmission(submission *Submission) (err error) {
	err = DB.Create(&submission).Error
	return
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

// GetProblemAllLastestSubmissionID 拿這堂課所有人最後一個的Problemsubmission
func GetProblemAllLastestSubmissionID(id uint) (submissions []string, err error) {
	var problem Problem
	var classuser []ClassUser
	if err = DB.Where("id = ?", id).First(&problem).Error; err != nil {
		return
	}
	if err = DB.Where(&ClassUser{Class_ID: problem.Class_ID}).Find(&classuser).Error; err != nil {
		return
	}

	for _, user := range classuser {
		var submission Submission
		if err = DB.Where(&Submission{UserID: user.User_ID, ProblemID: id}).Order("created_at desc").First(&submission).Error; err != nil {
			continue
		} else {
			submissions = append(submissions, strconv.Itoa(int(submission.PrivateSubmissionID)))
		}

	}

	return submissions, nil
}

// GetOwnAllSubmission 學生列出自己的submission
func GetOwnAllSubmission(id uint, user_id uint) (submissions []string, err error) {
	var problem Problem
	var classuser []ClassUser
	var submission Submission

	if err = DB.Where("id = ?", id).First(&problem).Error; err != nil {
		return
	}
	if err = DB.Where(&ClassUser{Class_ID: problem.Class_ID}).Find(&classuser).Error; err != nil {
		return
	}
	if err = DB.Where(&Submission{UserID: user_id, ProblemID: id}).Order("created_at desc").First(&submission).Error; err != nil {
		return []string{""}, nil
	} else {
		submissions = append(submissions, strconv.Itoa(int(submission.PrivateSubmissionID)))
	}

	return submissions, nil
}

//SubmissionBySubmissionID 用 submission ID 查詢submission
func SubmissionBySubmissionID(id uint) (Submission, error) {
	var submission Submission

	if err := DB.First(&submission, id).Error; err != nil {
		return Submission{}, err
	}

	return submission, nil
}

//ProblemByProblemID 用 作業 ID 查詢作業
func ProblemByProblemID(id uint) (Problem, error) {
	var problem Problem

	if err := DB.First(&problem, id).Error; err != nil {
		return Problem{}, err
	}

	return problem, nil
}
