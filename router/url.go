package router

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NCNUCodeOJ/BackendClassManagement/view"
	"github.com/gin-contrib/cors"

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

	// CORS
	if os.Getenv("FrontendURL") != "" {
		origins := strings.Split(os.Getenv("FrontendURL"), ",")
		log.Println("CORS:", origins)
		r.Use(cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
			AllowHeaders:     []string{"Origin, Authorization, Content-Type, Accept"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	r.GET("/ping", view.Pong)
	class := r.Group(baseURL + "/class")
	class.Use(authMiddleware.MiddlewareFunc())
	class.Use(getUserID())
	{
		// 課程
		class.GET("/:class_id", view.GetClassByID)  // 查詢課程
		class.POST("", view.CreateClass)            // 創建課程
		class.PATCH("/:class_id", view.UpdateClass) // 編輯課程資訊
		class.GET("", view.ListClass)               // 列出所有課程
		// 課程使用者
		class.GET("/:class_id/classuser/:classuser_id", view.GetClassUserByID)   // 查詢課程使用者
		class.POST("/:class_id/classuser", view.CreateClassUser)                 // 新增課程使用者
		class.DELETE("/:class_id/classuser/:classuser_id", view.DeleteClassUser) // 刪除課程使用者
		class.PATCH("/:class_id/classuser/:classuser_id", view.UpdateClassUser)  // 編輯課程使用者
		class.GET("/:class_id/classuser", view.ListClassUser)                    // 列出所有課程使用者

		// test 測驗
		class.GET("/:class_id/test/:test_id", view.GetTestByID)   // 查詢測驗
		class.POST("/:class_id/test", view.CreateTest)            // 創建測驗
		class.DELETE("/:class_id/test/:test_id", view.DeleteTest) // 刪除測驗
		class.PATCH("/:class_id/test/:test_id", view.UpdateTest)  // 編輯測驗資訊
		class.GET("/:class_id/test", view.ListTest)               // 列出所有課程測驗
	}
	// Problem 題目
	problem := r.Group(baseURL + "/class/:class_id/problem")
	problem.Use(authMiddleware.MiddlewareFunc())
	problem.Use(getUserID())
	{
		problem.POST("", view.CreateProblem)             // 創建程式碼題目
		problem.GET("/:problem_id", view.GetProblemByID) // 查詢程式碼題目資訊
		problem.GET("", view.ListProblem)                // 列出所有課程題目
		problem.DELETE("/:problem_id", view.DeleteProblem)
		problem.PATCH("/:problem_id", view.UpdateProblemQuestion)          // 編輯程式碼題目資訊
		problem.POST("/:problem_id/testcase", view.UploadQuestionTestCase) // 上傳程式碼題目測試 test case
	}
	questionsubmission := r.Group(baseURL + "/class/:class_id/problem/:problem_id")
	questionsubmission.Use(authMiddleware.MiddlewareFunc())
	questionsubmission.Use(getUserID())
	{
		questionsubmission.POST("/submission", view.CreateProblemSubmission)                // 上傳 submission
		questionsubmission.GET("/submission/:submission_id", view.GetProblemSubmissionByID) // 拿submission資訊
		questionsubmission.GET("/submission", view.ListSubmission)                          // 列出所有submission
	}
	mosssetup := r.Group(baseURL + "/class/:class_id/problem/:problem_id/moss") // 呼叫moss
	mosssetup.Use(authMiddleware.MiddlewareFunc())
	mosssetup.Use(getUserID())
	{
		mosssetup.GET("", view.SetupMoss) // 上傳 SetUp Moss
	}
	mosss := r.Group(privateURL + "/class/:class_id/problem/:problem_id/moss")
	mosss.Use(getUserID())
	{
		mosss.POST("", view.UploadMoss)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})
	return r
}
