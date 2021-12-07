package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/NCNUCodeOJ/BackendClassManagement/models"
	"github.com/NCNUCodeOJ/BackendClassManagement/mossservice"
	"github.com/NCNUCodeOJ/BackendClassManagement/router"
	"github.com/NCNUCodeOJ/BackendClassManagement/view"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjo0NzYwNjk2NDkyLCJpZCI6IjEiLCJvcmlnX2lhdCI6MTYzODYzMjQ5MiwidGVhY2hlciI6dHJ1ZSwidXNlcm5hbWUiOiJ2aW5jZW50In0.SUnwDQX_wkWlZdTMyCjhqIX4TIIzYrrY7lTiR_E2K8tvQBU1pyUgja60K0xcF1_x0m-egvRJQmhix5l6wdoR6g"
var class1ID, problem1ID, problem2ID, problem3ID, test1ID, submission1ID, submission2ID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
	view.Setup()
	mossservice.Setup()
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
		"teacher": 1
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
		"user_id":        2,
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
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateClassUser(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"user_id": 2,
		"role": 1
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/2", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestDeleteClassUser(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser/2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestCreate(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"testpaper_id":        ` + strconv.Itoa(1) + `,
		"start_time":  "2021-12-03T15:04:05Z08:00",
		"end_time":   "2021-12-04T15:04:05Z08:00"
	}`)
	// 2021-12-03 17:10:00 +0800 UTC 時間格式
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件

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
func TestClassUserCreate2(t *testing.T) {
	var data = []byte(`{
		"class_id":       ` + strconv.Itoa(class1ID) + `,
		"user_id":        2,
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
func TestListClass(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	///:class_id/problem/:problem_id
	req, _ := http.NewRequest("GET", "/api/v1/class", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestListClassUser(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/classuser", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListTest(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/test", nil)
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

func TestProblemCreate(t *testing.T) {
	var data = []byte(`{
		"end_time": 1638692151,
		"start_time": 1638692151,
		"language": "python3",
		"problem_name":       "接龍遊戲2",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"program_name":	      "Main",
		"layer":              1,
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "456", "output": "789"}
		],
		"tags_list":          ["1"],
		"hastestcase": "True"
	}`)
	r := router.SetupRouter()
	res := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(res, req)
	body, _ := ioutil.ReadAll(res.Body)
	s := struct {
		ProblemID int    `json:"problem_id"`
		Message   string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	problem1ID = s.ProblemID

	assert.Equal(t, http.StatusCreated, res.Code)
}
func TestGetProblem(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestEditProblem(t *testing.T) {
	var data = []byte(`{
		"problem_name":       "龍遊戲",
		"sample":             [
			{"input": "456", "output": "789"},
			{"input": "123", "output": "456"},
			{"input": "789", "output": "123"}
		],
		"tags_list":          ["4"]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	data = []byte(`{
		"problem_name":       "龍遊戲",
		"sample":             [
			{"input": "789", "output": "123"}
		]
	}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ = http.NewRequest("PATCH", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetProblem2(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUploadQuestionTestCase(t *testing.T) {
	r := gofight.New()

	url := "/api/v1/class/" + strconv.Itoa(class1ID) + "/problem/" + strconv.Itoa(problem1ID) + "/testcase"
	r.POST(url).
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase2.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r = gofight.New()
	r.POST(url).
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
func TestGetProblem3(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test2ProblemCreate(t *testing.T) {
	var data = []byte(`{
		"end_time": 1670371704,
		"start_time": 1638692151,
		"language": "python3",
		"problem_name":       "接龍遊戲2",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"program_name":	      "Main",
		"layer":              1,
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "456", "output": "789"}
		],
		"tags_list":          ["1"]
	}`)
	r := router.SetupRouter()
	res := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(res, req)
	body, _ := ioutil.ReadAll(res.Body)
	s := struct {
		ProblemID int    `json:"problem_id"`
		Message   string `json:"message"`
	}{}

	json.Unmarshal(body, &s)
	problem2ID = s.ProblemID
	assert.Equal(t, http.StatusCreated, res.Code)
}
func TestUploadQuestionTestCase2(t *testing.T) {
	r := gofight.New()

	url := "/api/v1/class/" + strconv.Itoa(class1ID) + "/problem/" + strconv.Itoa(problem2ID) + "/testcase"
	r.POST(url).
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase2.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r = gofight.New()
	r.POST(url).
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
func TestQuestionSubmissionCreate2(t *testing.T) {
	r := gofight.New()

	r.POST("/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem2ID)+"/submission").
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
	r.POST("/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem2ID)+"/submission").
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
}
func TestGetProblemSubmission(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem2ID)+"/submission/"+strconv.Itoa(submission1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestListSubmission2(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem2ID)+"/submission", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func Test3ProblemCreate(t *testing.T) {
	var data = []byte(`{
		"end_time": 1670371704,
		"start_time": 1670371704,
		"language": "python3",
		"problem_name":       "接龍遊戲2",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"program_name":	      "Main",
		"layer":              1,
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "456", "output": "789"}
		],
		"tags_list":          ["1"]
	}`)
	r := router.SetupRouter()
	res := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(res, req)
	body, _ := ioutil.ReadAll(res.Body)
	s := struct {
		ProblemID int    `json:"problem_id"`
		Message   string `json:"message"`
	}{}

	json.Unmarshal(body, &s)
	problem3ID = s.ProblemID
	assert.Equal(t, http.StatusCreated, res.Code)
}
func TestListProblem(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMosss(t *testing.T) {
	var data = []byte(`{
		"url": "https://moss.tw.edu.cn/login"
	}`)
	r := router.SetupRouter()
	res := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID)+"/moss", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestMossSetup(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID)+"/moss", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProblem(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("DELETE", "/api/v1/class/"+strconv.Itoa(class1ID)+"/problem/"+strconv.Itoa(problem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCleanup(t *testing.T) {
	err := os.Remove("test.db")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
}
