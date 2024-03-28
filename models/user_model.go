package models

import (
	"bytes"
	"echochat/internals"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserModel struct{}

type UserRegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Verified        bool   `json:"verified"`
}

func RegisterUser(ur UserRegisterRequest) {
	jsonData, err := json.Marshal(ur)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := internals.C.Pb.Post("/api/collections/users/records", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	fmt.Println("Created a new user")
}
