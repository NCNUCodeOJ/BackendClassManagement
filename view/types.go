package view

type classAPIRequest struct {
	Class_Name *string `json:"class_name"`
}
type classuserAPIRequest struct {
	Class_ID *uint `json:"class_id"`
	User_ID  *uint `json:"user_id"`
	Role     *int  `json:"role"`
}

type testAPIRequest struct {
	Class_ID     *uint   `json:"class_id"`
	TestPaper_ID *uint   `json:"testpaper_id"`
	Start_time   *string `json:"start_time"`
	End_time     *string `json:"end_time"`
}
type questionAPIRequest struct {
	Message    string `json:"message"`
	Problem_ID uint   `json:"problem_id"`
}
type SampleTemplate struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type getproblemAPIRequest struct {
	ProblemID         uint             `json:"problem_id"` //  外面看到得problem id
	ProblemName       string           `json:"problem_name"`
	ClassID           uint             `json:"class_id"`
	Description       string           `json:"description"`
	InputDescription  string           `json:"input_description"`
	OutputDescription string           `json:"output_description"`
	MemoryLimit       uint             `json:"memory_limit"`
	CPUTime           uint             `json:"cpu_time"`
	Layer             uint8            `json:"layer"`
	Sample            []SampleTemplate `json:"samples"`
	TagsList          []string         `json:"tags_list"`
	Start_Time        uint             `json:"start_time"`
	End_Time          uint             `json:"end_time"`
}
type mossAPIRequest struct {
	Problem_ID    string `json:"problem_id"`
	Submission_ID string `json:"submission_id"`
	Language      string `json:"language"`
}
