package handlers

import (
	"compile-server/internal/compilation"
	"compile-server/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userCode models.Code
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&userCode); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	log.Printf("Received request for %s", userCode.Path)

	userFile, err := compilation.MakeFile(userCode.Path, userCode.Lang, userCode.UserName, userCode.TaskName)
	if err != nil || userFile == "" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch userCode.Lang {
	case "cpp":
		{
			err := compilation.MakeCPPfile(userCode.TaskName, userFile)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
	case "py":
		{
			err := compilation.MakePYfile(userCode.TaskName, userFile)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
	default:
		http.Error(w, fmt.Sprintf("Unsupported language: %s", userCode.Lang), http.StatusBadRequest)
	}

}
