package main

import (
	"aisu.ai/api/v2/cmd/server/shared/api"
	"aisu.ai/api/v2/internal/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateUser(context *gin.Context) {
	var userReq *user.User
	if err := context.BindJSON(userReq); err != nil {
		log.Printf("Failed to parse new user request")
		context.IndentedJSON(http.StatusBadRequest, api.NewApiError(string(api.InvalidRequest), err.Error()))
	}

	repository := user.NewUserRepository(dbClient)
	if err := repository.Save(userReq); err != nil {
		log.Printf("An error occurred while saving the user to the database")
		context.IndentedJSON(http.StatusInternalServerError, NewApiError(internalError, err.Error()))
	}

	context.IndentedJSON(http.StatusCreated, userReq)
}
