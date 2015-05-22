package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var destinyUserNotFound = errors.New("User not found")
var bungieSDPURL = "http://www.bungie.net/Platform/Destiny/SearchDestinyPlayer/1/%s/"

type bungieSDP struct {
	Response []struct {
		ID string `json:"membershipId"`
	}
	ErrorStatus string
	Message     string
}

func resolveDestinyId(xb1gt string) (string, error) {
	var responseDocument bungieSDP

	log.Println("looking up destiny ID for", xb1gt)

	resp, err := http.Get(fmt.Sprintf(bungieSDPURL, xb1gt))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var decoder = json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseDocument); err != nil {
		return "", err
	}
	if responseDocument.ErrorStatus != "Success" {
		errorMessage := fmt.Errorf(
			"GT: %s, Status: %s, Message: %s",
			xb1gt,
			responseDocument.ErrorStatus,
			responseDocument.Message,
		)
		return "", errorMessage
	}
	if responseDocument.Response == nil {
		return "", destinyUserNotFound
	}
	if len(responseDocument.Response) == 0 {
		return "", destinyUserNotFound
	}
	return responseDocument.Response[0].ID, nil
}
