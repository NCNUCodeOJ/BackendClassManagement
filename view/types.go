package view

type classAPIRequest struct {
	Class_Name *string `json:"class_name"`
	Teacher    *uint   `json:"teacher"`
}
type classuserAPIRequest struct {
	Class_ID *uint `json:"class_id"`
	User_ID  *uint `json:"user_id"`
	Role     *int  `json:"role"`
}
type problemAPIRequest struct {
	Class_ID   *uint   `json:"class_id"`
	Problem_ID *uint   `json:"problem_id"`
	Start_time *string `json:"start_time"`
	End_time   *string `json:"end_time"`
}
type testAPIRequest struct {
	Class_ID     *uint   `json:"class_id"`
	TestPaper_ID *uint   `json:"testpaper_id"`
	Start_time   *string `json:"start_time"`
	End_time     *string `json:"end_time"`
}
