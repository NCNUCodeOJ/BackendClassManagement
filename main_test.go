package main

import (
	"NCNUOJBackend/ClassManagement/models"
	"NCNUOJBackend/ClassManagement/router"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDc5MTA4MjEyMywiaWQiOiI3MTI0MTMxNTQxOTcxMTA3ODYiLCJvcmlnX2lhdCI6MTYzNzQ4MjEyMywidXNlcm5hbWUiOiJ0ZXN0X3VzZXIifQ.pznOSok8X7qv6FSIihJnma_zEy70TerzOs0QDZOq_4RPYOKSEOOYTZ9-VLm2P9XRldS17-7QrLFwjjfXyCodtA"
var class1ID, classproblem1ID, problem1ID, test1ID, submission1ID, submission2ID, submission3ID, submission4ID, submission5ID, submission6ID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
}

func TestPing(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestClassCreate(t *testing.T) {
	var data = []byte(`{
		"class_name":       "程設1",
		"teacher":        1

	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		ClassID int    `json:"class_id"`
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	class1ID = s.ClassID
	assert.Equal(t, http.StatusCreated, w.Code)
}
func TestGetClassByID(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUpdateClass(t *testing.T) {
	var data = []byte(`{
		"class_name":       "課程2"
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	data = []byte(`{
		"class_name":       "課程3"
	}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ = http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestClassUserCreate(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"user_id":        1,
		"role":  0
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
func TestGetClassUserByID(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateClassUser(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"user_id": 1,
		"role": 1
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/1", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestDeleteClassUser(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestProblemCreate(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"problem_id":        1,
		"start_time":  "2021-12-03T15:04:05Z08:00",
		"end_time":   "2021-12-04T15:04:05Z08:00"
	}`)
	// 2021-12-03 17:10:00 +0800 UTC 時間格式
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)

	s := struct {
		Message   string `json:"message"`
		ProblemID int    `json:"problem_id"`
		Token     string `json:"token"`
	}{}
	json.Unmarshal(body, &s)
	classproblem1ID = s.ProblemID
	assert.Equal(t, http.StatusCreated, w.Code)
}
func TestGetProblemByID(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	///:class_id/problem/:problem_id
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(classproblem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUpdateProblem(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"problem_id": 2
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(classproblem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestDeleteProblem(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(classproblem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestTestCreate(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"testpaper_id":        1,
		"start_time":  "2021-12-03T15:04:05Z08:00",
		"end_time":   "2021-12-04T15:04:05Z08:00"
	}`)
	// 2021-12-03 17:10:00 +0800 UTC 時間格式
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	///class/:class_id/test
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/test", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)

	s := struct {
		Message string `json:"message"`
		TestID  int    `json:"test_id"`
	}{}
	json.Unmarshal(body, &s)
	test1ID = s.TestID
	assert.Equal(t, http.StatusCreated, w.Code)
}
func TestGetTestByID(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	///:class_id/problem/:problem_id
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/test/"+strconv.Itoa(test1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUpdateTest(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"testpaper_id": 2
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/test/"+strconv.Itoa(test1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestDeleteTest(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/class/"+strconv.Itoa(class1ID)+"/test/"+strconv.Itoa(test1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestCreateSubmission(t *testing.T) {
	r := gofight.New()
	///class/:class_id/problem/:problem_id/problem
	problem1ID = 1
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission1ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission2ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission3ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission4ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission5ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/class/:class_id/problem/:problem_id/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission6ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
