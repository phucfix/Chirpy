package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	cases := []struct{
		name 		string
		userID		uuid.UUID 
		tokenSecret string
		expiresIn  	time.Duration
		expectedErr bool
	} {
		{
			userID: 		uuid.New(),
			tokenSecret:	"mysecretKey",
			expiresIn: 		1 * time.Hour,
			expectedErr: 	false,
		},
		{
			userID:			uuid.New(),
			tokenSecret: 	"oiweh923842",
			expiresIn: 		1 * time.Hour,
			expectedErr:    false,
		},
		{
			userID: 		uuid.New(),
			tokenSecret:	"mySecretKey",
			expiresIn: 		10 * time.Hour,
			expectedErr: 	false,
		},
	}
	for _, c := range cases {
		t.Run("Test", func(t *testing.T) {
			tokenString, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
			if (err != nil) != c.expectedErr {
				t.Errorf("MakeJWT() err = %v, expectedErr %v", err, c.expectedErr)
				return
			}

			if tokenString == "" {
				t.Errorf("Expected a non-empty JWT string but got empty")
				return
			}		
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "MyTokenSecret"
	expiredAt := 30 * time.Second
	tokenString, err := MakeJWT(userID, tokenSecret, expiredAt)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	cases := []struct {
		name 		string
		tokenString string
		tokenSecret string
		expectedErr bool
		timeout 	time.Duration
	} {
		{
			name: "Valid token string and token secret",
			tokenString: tokenString,
			tokenSecret: tokenSecret,
			expectedErr: false,
			timeout: 1 * time.Second,
		},
		{
			name: "Invalid token string",
			tokenString: "Helloworld",
			tokenSecret: tokenSecret,
			expectedErr: true,
			timeout: 1 * time.Second,
		},
		{
			name: "Invalid secret string",
			tokenString: tokenString,
			tokenSecret: "Huhuhu",
			expectedErr: true,
			timeout: 1 * time.Second,
		},
		{
			name: "Invalid secret string and token",
			tokenString: "Helloworld",
			tokenSecret: "Hiworld",
			expectedErr: true,
			timeout: 1 * time.Second,
		},
		{
			name: "Reject expired token",
			tokenString: tokenString,
			tokenSecret: tokenSecret,
			expectedErr: true,
			timeout: 45 * time.Second,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			time.Sleep(c.timeout)
			actualUserID, err := ValidateJWT(c.tokenString, c.tokenSecret)
			if (err != nil) != c.expectedErr {
				t.Errorf("ValidateJWT() error = %v, is Err %v", err, c.expectedErr)
				return
			}

			if err == nil {
				if actualUserID != userID {
					t.Errorf("Expected userId = %v but actual = %v", userID, actualUserID)
				}
			}
		})
	}
}
