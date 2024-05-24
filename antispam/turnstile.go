package antispam

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var api = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type body struct {
	Secret string `json:"secret"`
	Token  string `json:"response"`
}

type cfResponse struct {
	Timestamp  string   `json:"challenge_ts"`
	Hostname   string   `json:"hostname"`
	Action     string   `json:"action"`
	Cdata      string   `json:"cdata"`
	ErrorCodes []string `json:"error-codes"`
	Success    bool     `json:"success"`
}

/*
Turnstile reports whether a token is valid using the Turnstile API.

  - secret: secret key associated with the site key
  - token: the token passed from the Turnstile widget

https://developers.cloudflare.com/turnstile/get-started/server-side-validation/
*/
func Turnstile(secret string, token string) (bool, error) {
	// create the request body
	b := body{
		Secret: secret,
		Token:  token,
	}
	bJSON, _ := json.Marshal(b)

	// post the request
	resp, err := http.Post(api, "application/json", bytes.NewBuffer(bJSON))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, errors.New("HTTP " + strconv.Itoa(resp.StatusCode))
	}

	// parse the response
	var respData cfResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return false, err
	}

	// handle validation failure
	if !respData.Success {
		if len(respData.ErrorCodes) > 0 {
			errorMsg := strings.Join(respData.ErrorCodes, ", ")
			return false, errors.New(errorMsg)
		}
		return false, errors.New("validation failed")
	}
	// only return true if the success field is true
	if respData.Success {
		return true, nil
	}
	return false, errors.New("unhandled error")
}
