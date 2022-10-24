package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Passowrd string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	log.Println("hit here with1", err, requestPayload)

	if err != nil {

		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validater user against db

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {

		app.errorJSON(w, errors.New("invalid crendentials"), http.StatusBadRequest)
		return
	}
	log.Println("hit here with2", user, err)

	valid, err := user.PasswordMatches(requestPayload.Passowrd)
	log.Println("hit here with", valid, err)
	if err != nil || !valid {

		app.errorJSON(w, errors.New("invalid crendentials"), http.StatusBadRequest)
		return
	}
	// log auth
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// log authentication
	paylod := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJSON(w, http.StatusAccepted, paylod)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
