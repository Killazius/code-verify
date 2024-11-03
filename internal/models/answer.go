package models

type Answer struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

const (
	Build   = "build"
	Compile = "compile"
	Test    = "test"
	Success = "success"
	Timeout = "timeout"
)
