package router

import (
	"NCNUOJBackend/ClassManagement/view"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/views"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
		if err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "系統錯誤",
				"error":   err.Error(),
			})
		} else {
			c.Set("userID", uint(id))
			c.Next()
		}
	}
}

// SetupRouter index
func SetupRouter() *gin.Engine {
	if gin.Mode() == "test" {
		err := godotenv.Load(".env.test")
		if err != nil {
			log.Println("Error loading .env file")
		}
	} else if gin.Mode() == "debug" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "NCNUOJ",
		SigningAlgorithm: "HS512",
		Key:              []byte(os.Getenv("SECRET_KEY")),
		MaxRefresh:       time.Hour,
		TimeFunc:         time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	baseURL := "api/v1"
	privateURL := "api/private/v1"
	r := gin.Default()
	r.GET("/ping", view.Pong)
	class := r.Group(baseURL + "/class")
	class.Use(authMiddleware.MiddlewareFunc())
	class.Use(getUserID())
	{
		// 課程
		class.GET("/:class_id", view.GetClassByID)  // 查詢課程
		class.POST("", view.CreateClass)            // 創建課程
		class.PATCH("/:class_id", view.UpdateClass) // 編輯課程資訊
		// 課程使用者
		class.GET("/:class_id/classuser/:classuser_id", view.GetClassUserByID)   // 查詢課程使用者
		class.POST("/:class_id/classuser", view.CreateClassUser)                 // 新增課程使用者
		class.DELETE("/:class_id/classuser/:classuser_id", view.DeleteClassUser) // 刪除課程使用者
		class.PATCH("/:class_id/classuser/:classuser_id", view.UpdateClassUser)  // 編輯課程使用者
		// Problem 題目
		class.GET("/:class_id/problem/:problem_id", view.GetProblemByID)   // 查詢題目
		class.POST("/:class_id/problem", view.CreateProblem)               // 創建題目
		class.DELETE("/:class_id/problem/:problem_id", view.DeleteProblem) // 刪除題目
		class.PATCH("/:class_id/problem/:problem_id", view.UpdateProblem)  // 編輯題目資訊
		// test 測驗
		class.GET("/:class_id/test/:test_id", view.GetTestByID)   // 查詢測驗
		class.POST("/:class_id/test", view.CreateTest)            // 創建測驗
		class.DELETE("/:class_id/test/:test_id", view.DeleteTest) // 刪除測驗
		class.PATCH("/:class_id/test/:test_id", view.UpdateTest)  // 編輯測驗資訊
	}
	privateProblem := r.Group(privateURL + "/class/:class_id/problem/:problem_id/problem")
	privateProblem.Use(authMiddleware.MiddlewareFunc())
	privateProblem.Use(getUserID())
	{
		privateProblem.POST("/:id/submission", views.CreateSubmission) // 上傳 submission
	}
	// submission
	submission := r.Group(privateURL + "/class/:class_id/problem/:problem_id//submission")
	{
		submission.PATCH("/:id/judge", views.UpdateSubmissionJudgeResult) // 更新 submission judge result
		submission.PATCH("/:id/style", views.UpdateSubmissionStyleResult) // 更新 submission style result
		submission.GET("/:id", views.GetSubmissionByID)                   // 取得 submission
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}
