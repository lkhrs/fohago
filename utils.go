package main

import "net/http"

// GetField returns the value of a field from a form submission
func GetField(r *http.Request, field string) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	return r.Form.Get(field), nil
}
