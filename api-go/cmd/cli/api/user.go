package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type newUserRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type user struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func CreateUser(name string) (string, error) {
	request := &newUserRequest{Name: name, Gender: "male"}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(ApiBaseUrl+"/users", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed to create new user")
	}

	var user user
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(respBytes, &user); err != nil {
		return "", err
	}
	return user.Id, nil
}
