package auth

import (
	"net/http"
	"testing"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		name  		string
		pw    		string
		expectedErr bool
	} {
		{
			name:  "Valid password",
			pw:    "password123",
			expectedErr: false,
		},
		{
			name: "Empty password",
			pw:   "",
			expectedErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := HashPassword(c.pw)
			if (err != nil) != c.expectedErr {
				t.Errorf("HashPassword() err = %v, isErr %v", err, c.expectedErr)
				return
			}
			if err == nil {
				if len(got) == 0 {
					t.Errorf("HashPass() returned empty hash for pw: %v", c.pw)
				}
				// Indirectly test by checking if CheckPwHash works with it.
				checkErr := CheckPasswordHash(got, c.pw)
				if checkErr != nil {
					t.Errorf("HashPass() generated a hash that CheckPwHash() rejected: %v", checkErr)
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securedPassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	cases := []struct {
		name 		string
		pw 			string
		hash    	string
		expectedErr bool
	} {
		{
			name: 	"valid passwd and hash",
			pw: 	password,
			hash: 	hash,
			expectedErr: 	false,
		},
		{
			name: 	"invalid passwd ",
			pw: 	"invalidPassword",
			hash: 	hash,
			expectedErr: 	true,
		},
		{
			name:	"invalid hash",
			pw: 	password,
			hash: 	"invalidHash",
			expectedErr: 	true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := CheckPasswordHash(c.hash, c.pw)
			if (err != nil) != c.expectedErr {
				t.Errorf("CheckPasswordHash() error = %v, isErr %v", err, c.expectedErr)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name 			string
		header 			http.Header
		expectedToken 	string
		expectedErr 	bool
	} {
		{
			name: "Valid bearer token",
			header: http.Header{"Authorization": []string{"Bearer mytoken"}},
			expectedToken: "mytoken",
			expectedErr: false,
		},
		{
			name: "Empty bearer token",
			header: http.Header{"Authorization": []string{}},
			expectedToken: "",
			expectedErr: true,
		},
		{
			name: "Missing Authorization header",
			header: http.Header{},
			expectedToken: "",
			expectedErr: true,
		},
		{
			name: "Invalid value",
			header: http.Header{"Authorization": []string{"Hello mytoken"}},
			expectedToken: "",
			expectedErr: true,
		},
		{
			name: "Extra space value",
			header: http.Header{"Authorization": []string{" Bearer mytoken"}},
			expectedToken : "mytoken",
			expectedErr: false,
		},
		{
			name: "Extra space value 2",
			header: http.Header{"Authorization": []string{"Bearer   mytoken"}},
			expectedToken: "mytoken",
			expectedErr: false,
		},
		{
			name: "Bearer with tab",
			header: http.Header{"Authorization": []string{"Bearer\tmytoken"}},
			expectedToken: "mytoken",
			expectedErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := GetBearerToken(c.header)
			if (err != nil) != c.expectedErr {
				t.Errorf("GetBearerToken() error = %v, expect Err %v", err, c.expectedErr)
				return
			}

			if (actual != c.expectedToken) {
				t.Errorf("Expected token: %s, but actual is %s", c.expectedToken, actual)
				return
			}
		})
	}
}
