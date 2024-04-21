package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sendigi-server/dtos"
)

func ConvertGoogleToken(token string) (*dtos.GooglePayload, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", token))
	if err != nil {
		return nil, err
	}

	// destroy when done
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}

	// if the body of the json has "error",
	// skip and return it as error.
	if body["error"] != nil {
		return nil, errors.New("invalid token")
	}

	// bind JSON into GooglePayload struct
	var payload dtos.GooglePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
