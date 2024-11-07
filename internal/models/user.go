package models

type UserJson struct {
	Code     string `json:"code"`
	Lang     string `json:"lang"`
	TaskName string `json:"task_name"`
	Token    string `json:"token"`
}

type Response struct {
	Username string `json:"username"`
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
