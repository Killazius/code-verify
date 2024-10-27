package models

import "fmt"

type CompilationError struct {
	Msg    string
	Reason error
}

func (e *CompilationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Reason)
}

func HandleCommonError(err error) error {
	if err != nil {
		return &CompilationError{
			Msg:    "Во время работы с файлом пользователя произошла ошибка",
			Reason: err,
		}
	}
	return nil
}
