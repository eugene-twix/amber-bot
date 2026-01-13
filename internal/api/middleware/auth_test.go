// internal/api/middleware/auth_test.go
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestParseUserJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  int64
		wantErr bool
	}{
		{
			name:    "valid user",
			input:   `{"id":123456789,"first_name":"John","last_name":"Doe","username":"johndoe","language_code":"en"}`,
			wantID:  123456789,
			wantErr: false,
		},
		{
			name:    "minimal user",
			input:   `{"id":42}`,
			wantID:  42,
			wantErr: false,
		},
		{
			name:    "missing id",
			input:   `{"first_name":"John"}`,
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "empty object",
			input:   `{}`,
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := parseUserJSON(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if user.ID != tt.wantID {
				t.Errorf("got ID %d, want %d", user.ID, tt.wantID)
			}
		})
	}
}

func TestSplitJSONPairs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple pairs",
			input: `"id":123,"name":"test"`,
			want:  2,
		},
		{
			name:  "nested object",
			input: `"id":123,"data":{"key":"value"},"name":"test"`,
			want:  3,
		},
		{
			name:  "string with comma",
			input: `"id":123,"name":"test, with comma"`,
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs := splitJSONPairs(tt.input)
			if len(pairs) != tt.want {
				t.Errorf("got %d pairs, want %d: %v", len(pairs), tt.want, pairs)
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Compute secret key the same way as middleware
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	// Create valid initData
	authDate := time.Now().Unix()
	userJSON := `{"id":123456789,"first_name":"John","username":"johndoe"}`

	// Build data-check-string
	params := map[string]string{
		"auth_date": fmt.Sprintf("%d", authDate),
		"user":      userJSON,
		"query_id":  "AAHdF6IQAAAAAN0XohDhrOrc",
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	dataCheckString := strings.Join(parts, "\n")

	// Compute hash
	hm := hmac.New(sha256.New, secretKey)
	hm.Write([]byte(dataCheckString))
	validHash := hex.EncodeToString(hm.Sum(nil))

	// Build raw initData
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	values.Set("hash", validHash)
	rawInitData := values.Encode()

	// Create middleware
	mw := &AuthMiddleware{
		botToken:  botToken,
		secretKey: secretKey,
	}

	t.Run("valid hash", func(t *testing.T) {
		if !mw.validateHash(rawInitData, validHash) {
			t.Error("expected valid hash to pass")
		}
	})

	t.Run("invalid hash", func(t *testing.T) {
		if mw.validateHash(rawInitData, "invalidhash123") {
			t.Error("expected invalid hash to fail")
		}
	})

	t.Run("tampered data", func(t *testing.T) {
		// Change auth_date in raw data but keep old hash
		tamperedValues := url.Values{}
		for k, v := range params {
			tamperedValues.Set(k, v)
		}
		tamperedValues.Set("auth_date", fmt.Sprintf("%d", authDate+100))
		tamperedValues.Set("hash", validHash)
		tamperedRaw := tamperedValues.Encode()

		if mw.validateHash(tamperedRaw, validHash) {
			t.Error("expected tampered data to fail validation")
		}
	})
}

func TestParseInitData(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	mw := &AuthMiddleware{
		botToken:  botToken,
		secretKey: secretKey,
	}

	t.Run("valid initData", func(t *testing.T) {
		authDate := time.Now().Unix()
		userJSON := `{"id":123,"first_name":"Test"}`

		values := url.Values{}
		values.Set("auth_date", fmt.Sprintf("%d", authDate))
		values.Set("user", userJSON)
		values.Set("query_id", "test123")
		values.Set("hash", "somehash")

		raw := values.Encode()

		initData, err := mw.parseInitData(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if initData.User.ID != 123 {
			t.Errorf("got user ID %d, want 123", initData.User.ID)
		}
		if initData.QueryID != "test123" {
			t.Errorf("got query_id %s, want test123", initData.QueryID)
		}
	})

	t.Run("missing hash", func(t *testing.T) {
		values := url.Values{}
		values.Set("auth_date", fmt.Sprintf("%d", time.Now().Unix()))
		values.Set("user", `{"id":123}`)

		_, err := mw.parseInitData(values.Encode())
		if err == nil {
			t.Error("expected error for missing hash")
		}
	})

	t.Run("missing auth_date", func(t *testing.T) {
		values := url.Values{}
		values.Set("user", `{"id":123}`)
		values.Set("hash", "somehash")

		_, err := mw.parseInitData(values.Encode())
		if err == nil {
			t.Error("expected error for missing auth_date")
		}
	})

	t.Run("missing user", func(t *testing.T) {
		values := url.Values{}
		values.Set("auth_date", fmt.Sprintf("%d", time.Now().Unix()))
		values.Set("hash", "somehash")

		_, err := mw.parseInitData(values.Encode())
		if err == nil {
			t.Error("expected error for missing user")
		}
	})
}

func TestAuthDateTTL(t *testing.T) {
	tests := []struct {
		name      string
		authDate  time.Time
		ttl       time.Duration
		wantValid bool
	}{
		{
			name:      "fresh auth",
			authDate:  time.Now(),
			ttl:       ReadTTL,
			wantValid: true,
		},
		{
			name:      "within read TTL",
			authDate:  time.Now().Add(-30 * time.Minute),
			ttl:       ReadTTL,
			wantValid: true,
		},
		{
			name:      "expired read TTL",
			authDate:  time.Now().Add(-2 * time.Hour),
			ttl:       ReadTTL,
			wantValid: false,
		},
		{
			name:      "within write TTL",
			authDate:  time.Now().Add(-5 * time.Minute),
			ttl:       WriteTTL,
			wantValid: true,
		},
		{
			name:      "expired write TTL",
			authDate:  time.Now().Add(-15 * time.Minute),
			ttl:       WriteTTL,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := time.Since(tt.authDate) <= tt.ttl
			if isValid != tt.wantValid {
				t.Errorf("auth_date validity: got %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}
