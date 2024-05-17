package main

import (
	"aisu.ai/api/v2/cmd/server/shared/api"
	"aisu.ai/api/v2/internal/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NewUserRequest struct {
	Name   string       `json:"name"`
	Gender *user.Gender `json:"gender"`
}

func CreateUser(context *gin.Context) {
	var request = &NewUserRequest{}
	if err := context.BindJSON(request); err != nil {
		errMsg := "An error occurred while parsing the request to create a new user"
		api.BindApiErrorResponse(context, errMsg, http.StatusBadRequest, api.ErrorCodeBadRequest, err)
		return
	}

	newUser := user.NewUser(request.Name, *request.Gender)
	createdUser, err := userRepository.Save(newUser)
	if err != nil {
		errMsg := "An error occurred while saving the user"
		api.BindApiErrorResponse(context, errMsg, http.StatusInternalServerError, api.ErrorCodeGeneralError, err)
	}

	context.IndentedJSON(http.StatusCreated, createdUser)

	//TODO: Bind the location header.
}

// func GetUser(context *gin.Context) {
//
// }
