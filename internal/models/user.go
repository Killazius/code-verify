package models

type Code struct {
	Path     string `json:"path"`
	Lang     string `json:"lang"`
	TaskName string `json:"task_name"`
	UserName string `json:"username"`
}

const (
	LangCpp = "cpp"
	LangPy  = "py"
)
