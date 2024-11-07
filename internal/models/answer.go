package models

type Answer struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}
type TokenAnswer struct {
	Status int `json:"status"`
}

const (
	Build   = "build"
	Compile = "compile"
	Test    = "test"
	Success = "success"
	Timeout = "timeout"
)
