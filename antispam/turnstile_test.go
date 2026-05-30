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

// https://developers.cloudflare.com/turnstile/troubleshooting/testing/
var keys = testKeys{
	Secret: testKey{
		Pass: "1x0000000000000000000000000000000AA",
		Fail: "2x0000000000000000000000000000000AA",
	},
	Token: "XXXX.DUMMY.TOKEN.XXXX",
}

func TestTurnstile(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		response    string
		expected    bool
		expectedErr error
	}{
		{
			name:     "valid token passes",
			secret:   keys.Secret.Pass,
			response: keys.Token,
			expected: true,
		},
		{
			name:        "invalid response returns typed error",
			secret:      keys.Secret.Fail,
			response:    keys.Token,
			expected:    false,
			expectedErr: ErrInvalidInputResponse,
		},
		{
			name:        "empty secret returns typed error",
			secret:      "",
			response:    keys.Token,
			expected:    false,
			expectedErr: ErrMissingInputSecret,
		},
		{
			name:        "empty response returns typed error",
			secret:      keys.Secret.Pass,
			response:    "",
			expected:    false,
			expectedErr: ErrMissingInputResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Turnstile(test.secret, test.response)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("Expected error %v, got %v", test.expectedErr, err)
			}
		})
	}
}
