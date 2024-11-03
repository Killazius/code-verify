package models

type UserJson struct {
	Code     string `json:"code"`
	Lang     string `json:"lang"`
	TaskName string `json:"task_name"`
	UserName string `json:"username"`
}

type Answer struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

const (
	LangCpp = "cpp"
	LangPy  = "py"
)

const (
	BaseCpp  = "base.cpp"
	BasePy   = "base.py"
	TestsTxt = "test.txt"
	TestCpp  = "test_cpp.go"
	TestPy   = "test_py.go"
)
