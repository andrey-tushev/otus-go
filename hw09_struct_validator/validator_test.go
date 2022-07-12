package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	BadTTag struct {
		Field string `json:"id" validate:"blah-blah-blah"`
	}

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

func TestNonValidationErrors(t *testing.T) {
	err := Validate(BadTTag{"something"})
	require.Error(t, err, ErrBadRule)

	err = Validate("not a struct type")
	require.Error(t, err, ErrNotAStruct)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{App{Version: "12345"}, nil},
		{App{Version: "1234500"}, ValidationError{Field: "Version", Err: ErrWrongLength}},
		{App{Version: "123"}, ValidationError{Field: "Version", Err: ErrWrongLength}},

		{Token{}, nil},
		{Token{Header: []byte("aaaaa"), Payload: []byte("bbbbb"), Signature: []byte("ccccc")}, nil},

		{Response{Code: 200}, nil},
		{Response{Code: 404}, nil},
		{Response{}, ValidationError{Field: "Code", Err: ErrIllegalValue}},
		{Response{Code: 100}, ValidationError{Field: "Code", Err: ErrIllegalValue}},

		{
			User{
				ID:     "012345678901234567890123456789ABCDEF",
				Name:   "Andrey",
				Age:    43,
				Email:  "andrey@example.com",
				Role:   "admin",
				Phones: []string{"71231234567", "71231234568", "71231234569"},
				meta:   nil,
			},
			nil,
		},
		{
			User{
				ID:     "-012345678901234567890123456789ABCDEF",
				Name:   "Andrey",
				Age:    143,
				Email:  "and#rey@example.com",
				Role:   "boss",
				Phones: []string{"71231234567", "71231234568", "123"},
				meta:   nil,
			},
			errors.New(
				"ID has wrong length, " +
					"Age is too big, " +
					"Email has bad format, " +
					"Role contains illegal value, " +
					"Phones has wrong length",
			),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &ValidationErrors{})
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
