package handlers

import (
	"compile-server/internal/compilation"
	"compile-server/internal/models"
	"encoding/json"
	"net/http"
)

func CodeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userCode models.Code
	err := json.NewDecoder(r.Body).Decode(&userCode)

	if err != nil {
		return
	}

	compilation.MakeFile(userCode.Link, userCode.Lang)
	if userCode.Lang == "cpp" {
		compilation.MakeCPPfile(userCode.Task_Name, "user.cpp")
	}
}
