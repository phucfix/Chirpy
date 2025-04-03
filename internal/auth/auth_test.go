package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		name  string
		pw    string
		isErr bool
	} {
		{
			name:  "Valid password",
			pw:    "password123",
			isErr: false,
		},
		{
			name: "Empty password",
			pw:   "",
			isErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := HashPassword(c.pw)
			if (err != nil) != c.isErr {
				t.Errorf("HashPassword() err = %v, isErr %v", err, c.isErr)
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
		name 	string
		pw 		string
		hash    string
		isErr 	bool
	} {
		{
			name: 	"valid passwd and hash",
			pw: 	password,
			hash: 	hash,
			isErr: 	false,
		},
		{
			name: 	"invalid passwd ",
			pw: 	"invalidPassword",
			hash: 	hash,
			isErr: 	true,
		},
		{
			name:	"invalid hash",
			pw: 	password,
			hash: 	"invalidHash",
			isErr: 	true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := CheckPasswordHash(c.hash, c.pw)
			if (err != nil) != c.isErr {
				t.Errorf("CheckPasswordHash() error = %v, isErr %v", err, c.isErr)
			}
		})
	}
}
