package antispam

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

var (
	ErrInvalidInputResponse = errors.New("invalid-input-response")
	ErrMissingInputSecret   = errors.New("missing-input-secret")
	ErrMissingInputResponse = errors.New("missing-input-response")
	ErrValidationFailed     = errors.New("validation failed")
)

// turnstileError maps Turnstile error codes to Go errors. If no codes are provided, it returns a generic validation failed error.
func turnstileError(codes []string) error {
	if len(codes) == 0 {
		return ErrValidationFailed
	}

	errs := make([]error, 0, len(codes))
	for _, code := range codes {
		switch code {
		case ErrInvalidInputResponse.Error():
			errs = append(errs, ErrInvalidInputResponse)
		case ErrMissingInputSecret.Error():
			errs = append(errs, ErrMissingInputSecret)
		case ErrMissingInputResponse.Error():
			errs = append(errs, ErrMissingInputResponse)
		default:
			errs = append(errs, errors.New(code))
		}
	}

	return errors.Join(errs...)
}

type HTTPStatusError struct {
	StatusCode int
}

func (e HTTPStatusError) Error() string {
	return "HTTP " + strconv.Itoa(e.StatusCode)
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

	// parse the response
	var respData cfResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 && len(respData.ErrorCodes) == 0 {
		return false, HTTPStatusError{StatusCode: resp.StatusCode}
	}

	// handle validation failure
	if !respData.Success {
		if len(respData.ErrorCodes) > 0 {
			return false, turnstileError(respData.ErrorCodes)
		}
		return false, ErrValidationFailed
	}
	// only return true if the success field is true
	if respData.Success {
		return true, nil
	}
	return false, errors.New("unhandled error")
}
