package antispam

import (
	"errors"
	"testing"
)

type testKeys struct {
	Secret testKey
	Token  string
}

type testKey struct {
	Pass string
	Fail string
}

var keys = testKeys{
	Secret: testKey{
		Pass: "1x0000000000000000000000000000000AA",
		Fail: "2x0000000000000000000000000000000AA",
	},
	Token: "XXXX.DUMMY.TOKEN.XXXX",
}

func TestTurnstile(t *testing.T) {
	// Test case 1: Pass
	secret := keys.Secret.Pass
	response := keys.Token
	expected := true
	result, err := Turnstile(secret, response)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 2: Fail
	secret = keys.Secret.Fail
	expected = false
	expectedErr := errors.New("invalid-input-response")
	result, err = Turnstile(secret, response)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error 'invalid-input-response', got %v", err)
	}
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 3: Empty secret
	secret = ""
	expectedErr = errors.New("missing-input-secret")
	_, err = Turnstile(secret, response)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error '%v', but got '%v'", expectedErr, err)
	}

	// Test case 4: Empty response
	secret = keys.Secret.Pass
	response = ""
	expectedErr = errors.New("missing-input-response")
	_, err = Turnstile(secret, response)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error '%v', but got '%v'", expectedErr, err)
	}
}
