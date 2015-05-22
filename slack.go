package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type slackUserListResponse struct {
	Ok      bool          `json:"ok"`
	Error   string        `json:"error"`
	Members []slackMember `json:"members"`
}

type slackMember struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Bot     bool   `json:"is_bot"`
	Profile struct {
		FirstName string `json:"first_name"`
		// Note: There's more stuff in the response,
		// but I don't need all that for this...
	} `json:"profile"`
	// Note: There's more stuff in the response,
	// but I don't need all that for this...
}

func getSlackUsers() ([]slackMember, error) {
	var responseDocument slackUserListResponse
	resp, err := http.Get(
		fmt.Sprintf(
			"https://slack.com/api/users.list?token=%s",
			url.QueryEscape(slackToken),
		),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var jsonDecoder = json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(&responseDocument); err != nil {
		return nil, err
	}
	if responseDocument.Ok {
		return responseDocument.Members, nil
	}
	return nil, errors.New(responseDocument.Error)
}
