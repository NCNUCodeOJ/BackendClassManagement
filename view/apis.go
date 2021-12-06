package view

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NCNUCodeOJ/BackendClassManagement/models"
	"github.com/buger/jsonparser"

	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

var hostURL = os.Getenv("HOST_URL")   // problem
var hostURL2 = os.Getenv("HOST_URL2") // test
var privateURL = "/api/private/v1"

// Role 0 學生 1 助教 2 老師
// class 課程
// problem 題目
// test 測驗

// 檢查課程使用者在課堂操作的權限
func Check_UserRole(user_id uint, class_id uint) (int, error) {
	if classuser, err := models.ClassUserByClassUserID(user_id, class_id); err != nil {
		return -1, err
	} else {
		return classuser.Role, nil
	}
}

// CreateClass 新增課程
func CreateClass(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	data := classAPIRequest{} // 接收創建課程資料的 struct 欄位: Class_Name(課程名稱)
	var class models.Class
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	replace.Replace(&class, &data) // 這邊的 &class 只有課程名稱
	class.Teacher = userID         // 設老師的 id 為 現在這個 user_id
	// 創建課程
	if err := models.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}

	// 創老師的課堂使用者資料，不然後面會沒權限執行其他動作
	var teacher models.ClassUser
	teacher.Class_ID = class.ID
	teacher.User_ID = userID
	teacher.Role = 2
	// 新增課堂使用者
	if err := models.CreateClassUser(&teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"class_id": class.ID,
		"message":  "create class complete",
	})

}

//CreateClassUser 新增課程使用者 (助教、老師可用)
func CreateClassUser(c *gin.Context) {
	data := classuserAPIRequest{} // 接收新增課程使用者資料的struct，欄位: Class_ID(課程ID)、User_ID(欲新增的課程使用者ID)、Role(角色)
	var classuser models.ClassUser
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 先抓 URL 裡面的課堂ID，準備確認操作權限
	classID := uint(class_id)
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	replace.Replace(&classuser, &data)
	classuser.Class_ID = classID
	// 新增課堂使用者
	if err := models.CreateClassUser(&classuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"class_id":     classID,
		"classuser_id": classuser.User_ID,
		"message":      "課程使用者新增成功",
	})

}

// CreateTest 新增測驗(助教、老師可用)
func CreateTest(c *gin.Context) {

	data := testAPIRequest{} // 接收新增課程測驗資料的struct
	var test models.Test

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	replace.Replace(&test, &data)
	//
	if err := models.CreateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"test_id":  test.ID,
		"message":  "測驗創建成功",
		"class_id": classID,
	})
}

// DeleteClass 刪除課程 (老師可用)
func DeleteClass(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 class_id 才知道要刪除哪個課程
	var class_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	userID := c.MustGet("userID").(uint)
	// 確認操作權限，限老師(2)可用
	if user_role, err := Check_UserRole(userID, class_ID); err != nil || user_role != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	if zero.IsZero(class_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	// 檢查是否有這堂課
	if _, err := models.ClassByClassID(class_ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "不存在此課程",
		})
		return
	}
	// 刪除課程
	if err := models.DeleteClass(class_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class_id": Id,
		"message":  "課程刪除成功",
	})
}

// DeleteClassUser 刪除課程使用者 (老師可用)
func DeleteClassUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id")) // 抓 URL 的 classuser_id，才知道要刪除哪個課堂使用者
	var classuser_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id，才知道是哪堂課
	classID := uint(class_id)
	// 確認操作權限，限老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	if zero.IsZero(classuser_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	// 檢查這位使用者是否存在在課堂
	if _, err := models.ClassUserByClassUserID(classuser_ID, classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "此課程無該使用者",
		})
		return
	}
	// 移除在這堂課的課堂使用者
	if err := models.DeleteClassUser(classuser_ID, classID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程使用者刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"classuser_id": Id,
		"message":      "課程使用者刪除成功",
	})

}

// DeleteProblem 刪除題目
func DeleteProblem(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 裡的 problem_id ，才知道要刪除哪個題目
	var prbolem_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	if zero.IsZero(prbolem_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	// 檢查課堂是否有這個題目
	if _, err := models.ProblemByProblemID(prbolem_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在此課程",
		})
		return
	}
	// 刪除題目
	if err := models.DeleteProblem(prbolem_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problem_id": Id,
		"message":    "課程刪除成功",
	})
}

// DeleteTest 刪除測驗 (助教、老師可用)
func DeleteTest(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	var test_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if zero.IsZero(test_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if _, err := models.TestByTestID(test_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在此課程",
		})
		return
	}
	if err := models.DeleteTest(test_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"test_id": Id,
		"message": "課程刪除成功",
	})
}

// UpdateClass 更新課程 (老師可用)
func UpdateClass(c *gin.Context) {
	class_id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id，才知道要改哪個課
	var class models.Class
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	classID := uint(class_id)
	// 確認操作權限，限老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	data := classAPIRequest{} // 目前只接收，class name ，只開放改課堂的名字

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	// 檢查是否有這堂課
	if class, err = models.ClassByClassID(classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此課程",
		})
		return
	}
	replace.Replace(&class, &data)
	class.ID = classID // model 看 id 去改資料
	// 更新課程的資訊(限課的名字)
	if err := models.UpdateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class_id": class_id,
		"message":  "課程更新成功",
	})
}

// UpdateClassUser 更新課程使用者 (老師可用)
func UpdateClassUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id")) // 抓 URL 的 classuser_id ，才知道要改哪位課堂使用者
	classuser_id := uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	// 確認操作權限，限老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	data := classuserAPIRequest{}
	var classuser models.ClassUser

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}

	replace.Replace(&classuser, &data)
	// model 看 class_id user_id 去改資料
	classuser.Class_ID = classID
	classuser.User_ID = classuser_id
	// 檢查這堂課是否有這個使用者
	if classuser, err := models.ClassUserByClassUserID(classuser.User_ID, classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"class":   classuser.Class_ID,
			"message": "此課程無該使用者",
		})
		return
	}
	// 更新課堂使用者的資訊
	if err := models.UpdateClassUser(&classuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"classuser_id": classuser_id,
		"message":      "課程使用者更新成功",
	})

}

// UpdateTest 更新測驗
func UpdateTest(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	data := testAPIRequest{}
	var test models.Test
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}

	replace.Replace(&test, &data)
	test.ID = uint(Id) // model 看 id 去改資料
	if _, err := models.TestByTestID(test.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "找不到該測驗",
		})
	}
	if err := models.UpdateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"test_id": Id,
		"message": "測驗更新成功",
	})
}

// GetClassUserByID 用課程使用者 ID 查詢課程使用者資訊 (學生、助教、老師可用)
func GetClassUserByID(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id")) // 抓 URL 的 classuser_id ，才知道要查哪位課堂使用者
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	var classuser_id uint = uint(Id)
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，學生只能查自己的資訊
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 && userID != classuser_id {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if zero.IsZero(Id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不能為零",
		})
		return
	}
	// 查課堂使用者的資訊
	if classuser, err := models.ClassUserByClassUserID(classuser_id, classID); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":     classuser.Class_ID,
			"classuser_id": classuser.User_ID,
			"role":         classuser.Role,
		})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"class":   classuser.Class_ID,
			"message": "此課程無使用者",
		})
		return
	}
}

// GetClassByID 用 ClassID 查詢課程 (學生、助教、老師可用)
func GetClassByID(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道要找哪堂課
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	var class_id uint = uint(Id)
	userID := c.MustGet("userID").(uint)
	// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，只能查自己有的課程
	if _, err := Check_UserRole(userID, class_id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if zero.IsZero(Id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不能為零",
		})
		return
	}
	// 查自己的課堂資訊
	if class, err := models.ClassByClassID(class_id); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":   class.ID,
			"class_name": class.Class_Name,
			"teacher":    class.Teacher,
		})
		return
	}
}

// GetTestByID 用測驗 id 查詢測驗 (輸出未完成)
func GetTestByID(c *gin.Context) {

	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	var test_id uint = uint(Id)
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Permission denied",
		})
		return
	}
	if zero.IsZero(Id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不能為零",
		})
		return
	}
	if test, err := models.TestByTestID(test_id); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":     test.Class_ID,
			"testPaper_id": test.TestPaper_ID,
			"start_time":   test.Start_time,
			"end_time":     test.End_time,
		})
		return
	}
}

// ListClass 列出使用者所有課堂 (學生、助教、老師可用)
func ListClass(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	// 只會列出自己有的課堂
	if classes, err := models.ListClassUserClass(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		var classes_id []uint
		for _, data := range classes {
			classes_id = append(classes_id, data.Class_ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"classes": classes_id,
			"message": "list class complete",
		})
		return
	}
}

// ListClassUser 列出所有課堂使用者 (學生、助教、老師可用)
func ListClassUser(c *gin.Context) {
	class_id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，使用者要在該課程，才能查列出所有該課的課程使用者
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	// 列出所有該課的課程使用者
	if classusers, err := models.ListClassUser(classID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		var classusers_list []uint
		for _, data := range classusers {
			classusers_list = append(classusers_list, data.User_ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"classusers": classusers_list,
			"message":    "list classuser complete",
		})
		return
	}
}

// ListProblem 列出課堂所有題目 (學生、助教、老師可用)
func ListProblem(c *gin.Context) {
	class_id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，使用者要在該課程，才能查列出所有該課的課程題目
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	// 列出課堂所有題目
	if problems, err := models.ListProblem(classID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		var problems_list []uint
		for _, data := range problems {
			problems_list = append(problems_list, data.ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"problems": problems_list,
			"message":  "list classproblem complete",
		})

		return
	}
}

// ListTest 列出課堂所有測驗
func ListTest(c *gin.Context) {
	class_id, err := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	if tests, err := models.ListTest(classID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		var tests_list []uint
		for _, data := range tests {
			tests_list = append(tests_list, data.ID)
		}
		c.JSON(http.StatusOK, gin.H{
			"tests":   tests_list,
			"message": "list classtest complete",
		})
		return
	}
}

// CreateProblem 創程式碼題目 (助教、老師可用)
func CreateProblem(c *gin.Context) {
	if gin.Mode() == "test" {
		var problem_test models.Problem
		rawdata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system1 error",
			})
			return
		}
		class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
		classID := uint(class_id)

		userID := c.MustGet("userID").(uint)
		// 確認操作權限，限助教(1)、老師(2)可用
		if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Permission denied",
			})
			return
		}
		// 抓取其中的開始時間
		if start_time, err := jsonparser.GetInt(rawdata, "start_time"); err == nil {
			problem_test.Start_time = time.Unix(start_time, 0)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "開始時間未填寫",
			})
		}
		// 抓取其中的結束時間
		if end_time, err := jsonparser.GetInt(rawdata, "end_time"); err == nil {
			problem_test.End_time = time.Unix(end_time, 0)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "結束時間未填寫",
			})
		}
		problem_test.Problem_ID = uint(123)
		// 新增該堂課的題目
		if err := models.CreateProblem(&problem_test); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "系統錯誤",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":    "題目創建成功",
			"problem_id": problem_test.ID,
		})
		return
	}
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)

	userID := c.MustGet("userID").(uint)
	question := questionAPIRequest{} // 獲取程式碼題目的 real ID
	var problem models.Problem
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	URL := hostURL + privateURL + "/problem"
	// 給 rawdata
	rawdata, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system1 error",
		})
		return
	}

	problem.Class_ID = classID               // 設 prbolem 的課堂ID
	responseBody := bytes.NewBuffer(rawdata) // 把 rawdata 塞進 body
	//Leverage Go's HTTP Post function to make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, responseBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 設 header 原封不動船過去
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req) // 進行請求
	//Handle Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	defer res.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 解析 Response 到 question
	if err := json.Unmarshal(body, &question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 確認是否創建成功
	if question.Problem_ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫 ",
		})
		return
	} else {
		problem.Problem_ID = question.Problem_ID
	}
	// 抓取其中的開始時間
	if start_time, err := jsonparser.GetInt(rawdata, "start_time"); err == nil {
		problem.Start_time = time.Unix(start_time, 0)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "開始時間未填寫",
		})
	}
	// 抓取其中的結束時間
	if end_time, err := jsonparser.GetInt(rawdata, "end_time"); err == nil {
		problem.End_time = time.Unix(end_time, 0)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "結束時間未填寫",
		})
	}
	// 新增該堂課的題目
	if err := models.CreateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "題目創建成功",
		"problem_id": problem.ID,
	})
}

// GetProblemByID 用題目ID查程式碼題目 (學生、助教、老師可用)
func GetProblemByID(c *gin.Context) {
	if gin.Mode() == "test" {
		class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
		classID := uint(class_id)

		userID := c.MustGet("userID").(uint)
		problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
		var problemID = uint(problem_id)
		// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，
		if _, err := Check_UserRole(userID, classID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Permission denied",
			})
			return
		}
		var problems models.Problem
		if problem, err := models.ProblemByProblemID(problemID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "problem doesn't exist",
			})
			return
		} else {
			problems = problem
		}
		c.JSON(http.StatusOK, gin.H{
			"message":            "problem exist",
			"problem_id":         problemID,
			"class_id":           classID,
			"start_time":         uint(problems.Start_time.Unix()),
			"end_time":           uint(problems.End_time.Unix()),
			"problem_name":       "接龍遊戲2",
			"description":        "開始接龍",
			"input_description":  "567",
			"output_description": "789",
			"memory_limit":       134217728,
			"cpu_time":           1000,
			"layer":              1,
			"samples":            `[{"input": "123", "output": "456"},{"input": "456", "output": "789"}]`,
			"tags_list":          `["簡單"]`,
		})

		return
	}
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)

	userID := c.MustGet("userID").(uint)
	problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
	var problemID = uint(problem_id)
	var problem_data getproblemAPIRequest // 接收回傳的題目資訊
	var question_ID int                   // real problem id
	// 確認操作權限，限學生(0)、助教(1)、老師(2)可用，
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	// 檢查該堂課是否有這個題目
	if problem, err := models.ProblemByProblemID(problemID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "problem doesn't exist",
		})
		return
	} else {
		question_ID = int(problem.Problem_ID)
	}

	URL := hostURL + privateURL + "/problem/" + strconv.Itoa(question_ID)

	//Leverage Go's HTTP Post function to make request
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})

		return
	}

	// 設置 header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	//Handle Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})

		return
	}
	defer res.Body.Close()
	//Read the response body

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})

		return
	}
	if err := json.Unmarshal(body, &problem_data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
	}
	// 確認回傳是否有東西
	if problem_data.ProblemName != "" {
		if problem, err := models.ProblemByProblemID(problemID); err == nil {
			problem_data.Start_Time = uint(problem.Start_time.Unix())
			problem_data.End_Time = uint(problem.End_time.Unix())
			problem_data.ProblemID = problemID
			problem_data.ClassID = problem.Class_ID
		}

		c.JSON(http.StatusOK, gin.H{
			"message":            "problem exist",
			"problem_id":         problem_data.ProblemID,
			"class_id":           problem_data.ClassID,
			"start_time":         problem_data.Start_Time,
			"end_time":           problem_data.End_Time,
			"problem_name":       problem_data.ProblemName,
			"description":        problem_data.Description,
			"input_description":  problem_data.InputDescription,
			"output_description": problem_data.OutputDescription,
			"memory_limit":       problem_data.MemoryLimit,
			"cpu_time":           problem_data.CPUTime,
			"layer":              problem_data.Layer,
			"samples":            problem_data.Sample,
			"tags_list":          problem_data.TagsList,
		})

		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})

		return
	}
}

// UpdateProblemQuestion 編輯question 或是時間 (助教、老師可用)
func UpdateProblemQuestion(c *gin.Context) {
	if gin.Mode() == "test" {
		class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
		classID := uint(class_id)

		userID := c.MustGet("userID").(uint)
		problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
		var problemID = uint(problem_id)
		rawdata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system1 error",
			})
			return
		}

		// 確認操作權限，限助教(1)、老師(2)可用
		if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Permission denied",
			})
			return
		}
		// 抓取其中的開始時間
		var problem models.Problem
		problem.ID = problemID
		problem.Problem_ID = uint(123)
		// 拿該題目的資料
		if data, err := models.ProblemByProblemID(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system error",
			})
			return
		} else {
			problem.Start_time = data.Start_time
			problem.End_time = data.End_time
		}
		if start_time, err := jsonparser.GetInt(rawdata, "start_time"); err == nil {
			problem.Start_time = time.Unix(start_time, 0)
		}
		// 抓取其中的結束時間
		if end_time, err := jsonparser.GetInt(rawdata, "end_time"); err == nil {
			problem.End_time = time.Unix(end_time, 0)
		}
		// 更新題目資訊
		if err := models.UpdateProblem(&problem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":    "題目編輯失敗",
				"problem_id": problemID,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message":    "題目編輯成功",
			"problem_id": problemID,
		})
		return
	}
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)

	userID := c.MustGet("userID").(uint)
	problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
	var problemID = uint(problem_id)
	var question_ID int
	// 確認操作權限，限助教(1)、老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	// 檢查是否有該題目
	if problem, err := models.ProblemByProblemID(problemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		question_ID = int(problem.Problem_ID)
	}

	URL := hostURL + privateURL + "/problem/" + strconv.Itoa(question_ID)

	// 要更改的資料
	rawdata, err := c.GetRawData()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	responseBody := bytes.NewBuffer(rawdata)
	//Leverage Go's HTTP Post function to make request
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", URL, responseBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
	}
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	//Handle Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	defer res.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	data := questionAPIRequest{}

	json.Unmarshal(body, &data)
	// 確認是否回傳成功
	// 抓取其中的開始時間
	var problem models.Problem
	problem.ID = problemID
	problem.Problem_ID = uint(question_ID)
	// 拿該題目的資料
	if data, err := models.ProblemByProblemID(problemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		problem.Start_time = data.Start_time
		problem.End_time = data.End_time
	}
	if start_time, err := jsonparser.GetInt(rawdata, "start_time"); err == nil {
		problem.Start_time = time.Unix(start_time, 0)
	}
	// 抓取其中的結束時間
	if end_time, err := jsonparser.GetInt(rawdata, "end_time"); err == nil {
		problem.End_time = time.Unix(end_time, 0)
	}
	// 更新題目資訊
	if err := models.UpdateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":    "題目編輯失敗",
			"problem_id": problemID,
		})
	}
	if data.Problem_ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "題目編輯成功",
			"problem_id": problemID,
		})
		return
	}

}

// UploadQuestionTestCase 上傳題目測試 testcase (老師可用)
func UploadQuestionTestCase(c *gin.Context) {
	if gin.Mode() == "test" {
		class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
		classID := uint(class_id)

		userID := c.MustGet("userID").(uint)
		problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
		var problemID = uint(problem_id)
		_, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system1 error",
			})
			return
		}

		// 確認操作權限，限助教(1)、老師(2)可用
		if user_role, err := Check_UserRole(userID, classID); err != nil || user_role < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Permission denied",
			})
			return
		}
		// 確認是否有題目
		if _, err := models.ProblemByProblemID(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system error",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":          "上傳成功",
			"problem_id":       problemID,
			"test_case_number": 1,
		})
		return
	}
	class_id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
	var problemID = uint(problem_id)

	var question_ID int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	userID := c.MustGet("userID").(uint)
	// 確認操作權限，限老師(2)可用
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}
	// 確認是否有題目
	if problem, err := models.ProblemByProblemID(problemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	} else {
		question_ID = int(problem.Problem_ID)
	}

	URL := hostURL + privateURL + "/problem/" + strconv.Itoa(question_ID) + "/testcase"
	rawdata, err := c.GetRawData() // 原始資料

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	responseBody := bytes.NewBuffer(rawdata)
	//Leverage Go's HTTP Post function to make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, responseBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 設 header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	res, err := client.Do(req)

	//Handle Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	defer res.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 回傳 api 回傳的東西
	c.Status(res.StatusCode)
	c.Writer.Write(body)
}

// CreateProblemSubmission 創建題目submission (學生、助教、老師可用)
func CreateProblemSubmission(c *gin.Context) {
	if gin.Mode() == "test" {
		class_id, _ := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
		classID := uint(class_id)

		userID := c.MustGet("userID").(uint)
		problem_id, _ := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
		var problemID = uint(problem_id)
		_, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system1 error",
			})
			return
		}

		// 確認操作權限，限限學生(0)、助教(1)、老師(2)可用
		if _, err := Check_UserRole(userID, classID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Permission denied",
			})
			return
		}
		// 確認是否有題目
		if _, err := models.ProblemByProblemID(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "system error",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":       "提交成功",
			"problem_id":    problemID,
			"submission_id": 000000000000000001,
		})
		return
	}
	class_id, err := strconv.Atoi(c.Params.ByName("class_id")) // 抓 URL 的 class_id ，才知道是哪堂課
	classID := uint(class_id)
	userID := c.MustGet("userID").(uint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	problem_ID, err := strconv.Atoi(c.Params.ByName("problem_id")) // 抓 URL 的 problem_id ，才知道是哪個題目
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	problemID := uint(problem_ID)
	// 確認操作權限，限學生(1)、助教(1)、老師(2)可用
	if _, err := Check_UserRole(userID, classID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	var question_id string
	// 檢查是否有該題目
	if data, err := models.ProblemByProblemID(problemID); err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Problem not found",
		})
		return
	} else {
		question_id = strconv.Itoa(int(data.Problem_ID))
	}
	URL := hostURL + privateURL + "/problem" + "/" + question_id + "/submission"

	rawdata, err := c.GetRawData() // 原始資料

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	responseBody := bytes.NewBuffer(rawdata)
	//Leverage Go's HTTP Post function to make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, responseBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 設 header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	//Handle Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	defer res.Body.Close()
	//Read the response body

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
		})
		return
	}
	// 回傳 api 回傳的東西
	c.Status(res.StatusCode)
	c.Writer.Write(body)

}

// CreateMoss 呼叫 moss API
