package models

import (
	"bytes"
	"echochat/internals"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UserResponse struct {
	Token  string `json:"token"`
	Record struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Roles string `json:"roles"`
	} `json:"record"`
}

type UserModel struct {
	Token string   `json:"token"`
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type UserRegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Verified        bool   `json:"verified"`
}

type UserLoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AuthError struct {
	Code    int8   `json:"code"`
	Message string `json:"message"`
}

type UserSession struct {
	CSRF string    `json:"csrf"`
	User UserModel `json:"user"`
}

func LoginUser(ul *UserLoginRequest) (UserModel, error) {
	jsonData, err := json.Marshal(ul)
	if err != nil {
		return UserModel{}, err
	}

	client := &http.Client{}
	resp, err := client.Post(fmt.Sprintf("%v%v", internals.C.BaseURL, "/api/collections/users/auth-with-password"),
		"application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return UserModel{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserModel{}, errors.New("failed response status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return UserModel{}, err
	}

	var userResponse UserResponse

	if err := json.Unmarshal(body, &userResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return UserModel{}, err
	}

	var roles []string

	rolesInput := strings.Trim(userResponse.Record.Roles, "[]")
	roleElements := strings.Split(rolesInput, ",")

	for _, element := range roleElements {
		element = strings.TrimSpace(element)
		element = strings.Trim(element, "'\"")

		roles = append(roles, element)
	}

	return UserModel{
		ID:    userResponse.Record.ID,
		Token: userResponse.Token,
		Email: userResponse.Record.Email,
		Roles: roles,
	}, nil
}
