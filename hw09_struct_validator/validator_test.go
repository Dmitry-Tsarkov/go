package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^[\\w.]+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid User",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "john.doe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543211"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid User",
			in: User{
				ID:    "123",
				Age:   17,
				Email: "john.doe@example",
				Role:  "guest",
				Phones: []string{
					"123456789",
					"09876543211",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("length of value is 3 but expected 36")},
				{Field: "Age", Err: errors.New("value 17 is less than minimum 18")},
				{Field: "Email", Err: errors.New("value john.doe@example does not match pattern ^[\\w.]+@\\w+\\.\\w+$")},
				{Field: "Role", Err: errors.New("value guest is not in the allowed list [admin,stuff]")},
				{Field: "Phones", Err: errors.New("length of value is 9 but expected 11")},
			},
		},
		{
			name: "valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid App",
			in: App{
				Version: "123",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: errors.New("length of value is 3 but expected 5")},
			},
		},
		{
			name: "valid Response",
			in: Response{
				Code: 200,
				Body: "Success",
			},
			expectedErr: nil,
		},
		{
			name: "invalid Response",
			in: Response{
				Code: 201,
				Body: "Success",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errors.New("value 201 is not in the allowed list [200,404,500]")},
			},
		},
		{
			name: "valid Token",
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if (err == nil) != (tt.expectedErr == nil) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error message %q, got %q", tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
