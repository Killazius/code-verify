package ws

import (
	"compile-server/internal/compilation"
	"encoding/json"
	"fmt"
)

type UserMessage struct {
	Code   string           `json:"code"`
	Lang   compilation.Lang `json:"lang"`
	TaskID string           `json:"task_id"`
	Token  string           `json:"token"`
}

func (u *UserMessage) UnmarshalJSON(data []byte) error {
	type Alias UserMessage
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if !compilation.IsValidLang(aux.Lang) {
		return fmt.Errorf("invalid language")
	}
	if u.TaskID == "" {
		return fmt.Errorf("invalid task id")
	}
	if u.Code == "" {
		return fmt.Errorf("invalid code")
	}
	return nil

}
