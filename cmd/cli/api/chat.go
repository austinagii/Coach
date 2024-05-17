package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type task struct {
	Objective string `json:"objective"`
}

type newChatRequest struct {
	UserId string `json:"user_id"`
	Task   *task  `json:"task"`
}

type chat struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func CreateChat(userId string) (string, string, error) {
	task := &task{Objective: "goal_creation"}
	request := &newChatRequest{UserId: userId, Task: task}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", "", err
	}
	resp, err := http.Post(ApiBaseUrl+"/chats", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", "", errors.New("failed to create new chat")
	}

	var chat chat
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	if err := json.Unmarshal(respBytes, &chat); err != nil {
		return "", "", err
	}
	return chat.Id, chat.Text, nil
}

type outgoingChatMessage struct {
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}

type incomingChatMessage struct {
	Text string `json:"text"`
}

func Respond(userId string, chatId string, userMessage string) (string, error) {
	request := &outgoingChatMessage{UserId: userId, Text: userMessage}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/chats/%s/messages", ApiBaseUrl, chatId)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get response to chat message")
	}

	var message incomingChatMessage
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(respBytes, &message); err != nil {
		return "", err
	}
	return message.Text, nil
}
