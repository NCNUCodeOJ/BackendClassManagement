package view

import (
	"NCNUOJBackend/ClassManagement/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

// Role 0 學生 1 助教 2 老師
// class 課程
// problem 題目
// test 測驗
// 檢查課程操作的權限
func Check_UserRole(user_id uint, class_id uint) (int, error) {
	if classuser, err := models.ClassUserByClassUserID(user_id, class_id); err != nil {
		return -1, err
	} else {
		return classuser.Role, nil
	}
}

// 新增課程
func CreateClass(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	data := classAPIRequest{}
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
	replace.Replace(&class, &data)
	class.Teacher = userID

	if err := models.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	// 一起創課堂使用者資料 不然沒權限

	var teacher models.ClassUser
	teacher.Class_ID = class.ID
	teacher.User_ID = userID
	teacher.Role = 2
	if err := models.CreateClassUser(&teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"class_id": class.ID,
		"message":  "課程創建成功",
	})

}

// 新增課程使用者 助教 老師
func CreateClassUser(c *gin.Context) {
	data := classuserAPIRequest{}
	var classuser models.ClassUser
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
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
	if err := models.CreateClassUser(&classuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"class_id":     classID,
		"classuser_id": classuser.User_ID,
		"message":      "課程使用者新增成功",
	})

}

// 新增題目 助教 老師可用
func CreateProblem(c *gin.Context) {
	token := c.GetHeader("A")
	data := problemAPIRequest{}
	var problem models.Problem

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
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
	replace.Replace(&problem, &data)
	if err := models.CreateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":    "題目創建成功",
		"problem_id": problem.ID,
		"class_id":   classID,
		"token":      token,
	})

}

// 新增測驗
func CreateTest(c *gin.Context) {
	data := testAPIRequest{}
	var test models.Test

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
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
	if err := models.CreateTest(&test); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"test_id":  test.ID,
		"message":  "測驗創建成功",
		"class_id": classID,
	})
}

// 刪除課程
func DeleteClass(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("class_id"))
	var class_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}

	userID := c.MustGet("userID").(uint)
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
	if _, err := models.ClassByClassID(class_ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "不存在此課程",
		})
		return
	}
	if err := models.DeleteClass(class_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "課程刪除成功",
	})
}

// 刪除課程使用者
func DeleteClassUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id"))
	var classuser_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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

	if zero.IsZero(classuser_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if _, err := models.ClassUserByClassUserID(classuser_ID, classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "此課程無該使用者",
		})
		return
	}
	if err := models.DeleteClassUser(classuser_ID, classID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程使用者刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "課程使用者刪除成功",
	})

}

// 刪除題目
func DeleteProblem(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("problem_id"))
	var prbolem_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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

	if zero.IsZero(prbolem_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if _, err := models.ProblemByProblemID(prbolem_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在此課程",
		})
		return
	}
	if err := models.DeleteProblem(prbolem_ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "課程刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "課程刪除成功",
	})
}

// 刪除測驗
func DeleteTest(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	var test_ID uint = uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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
	if zero.IsZero(test_ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if _, err := models.TestkByTestID(test_ID); err != nil {
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
		"message": "課程刪除成功",
	})
}

// 更新課程
func UpdateClass(c *gin.Context) {
	class_id, err := strconv.Atoi(c.Params.ByName("class_id"))
	var class models.Class
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	userID := c.MustGet("userID").(uint)
	classID := uint(class_id)
	if user_role, err := Check_UserRole(userID, classID); err != nil || user_role != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		return
	}

	data := classAPIRequest{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}
	if class, err = models.ClassByClassID(classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此課程",
		})
		return
	}
	replace.Replace(&class, &data)
	class.ID = classID // model 看 id 去改資料

	if err := models.UpdateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失敗",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "課程更新成功",
	})
}

// 更新課程使用者
func UpdateClassUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id"))
	classuser_id := uint(Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
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
	// 檢查是否有這個學生
	if classuser, err := models.ClassUserByClassUserID(classuser.User_ID, classID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"class":   classuser.Class_ID,
			"message": "此課程無該使用者",
		})
		return
	}
	if err := models.UpdateClassUser(&classuser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "課程使用者更新成功",
	})

}

// 更新題目
func UpdateProblem(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("problem_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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
	data := problemAPIRequest{}
	var problem models.Problem
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫",
		})
		return
	}

	replace.Replace(&problem, &data)
	problem.ID = uint(Id) // model 看 id 去改資料
	if _, err := models.ProblemByProblemID(problem.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "找不到該題目",
		})
	}
	if err := models.UpdateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "題目更新成功",
	})
}

// 更新測驗
func UpdateTest(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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
	if _, err := models.TestkByTestID(test.ID); err != nil {
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
		"message": "測驗更新成功",
	})
}

// 用課程使用者 ID 查詢課程使用者 (輸出未完成)
func GetClassUserByID(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("classuser_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}

	var classuser_id uint = uint(Id)
	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
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

// 用 ClassID 查詢課程 (輸出未完成)
func GetClassByID(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("class_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	var class_id uint = uint(Id)
	userID := c.MustGet("userID").(uint)

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
	if class, err := models.ClassByClassID(class_id); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":   class.ID,
			"class_name": class.Class_Name,
			"teacher":    class.Teacher,
		})
		return
	}
}

// 用題目 ID 查詢題目 (輸出未完成)
func GetProblemByID(c *gin.Context) {
	Id, err := strconv.Atoi(c.Params.ByName("problem_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	}
	var problem_id uint = uint(Id)

	userID := c.MustGet("userID").(uint)
	class_id, _ := strconv.Atoi(c.Params.ByName("class_id"))
	classID := uint(class_id)
	if _, err := Check_UserRole(userID, classID); err != nil {
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
	if problem, err := models.ProblemByProblemID(problem_id); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":   problem.Class_ID,
			"problem_id": problem.Problem_ID,
			"start_time": problem.Start_time,
			"end_time":   problem.End_time,
		})
		return
	}
}

// 用測驗 id 查詢測驗 (輸出未完成)
func GetTestByID(c *gin.Context) {
	c.GetHeader("A")
	Id, err := strconv.Atoi(c.Params.ByName("test_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
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
	if test, err := models.TestkByTestID(test_id); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"class_id":     test.Class_ID,
			"testPaper_id": test.TestPaper_ID,
			"start_time":   test.Start_time,
			"end_time":     test.End_time,
		})
		return
	}
}
