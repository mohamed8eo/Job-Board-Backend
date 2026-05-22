package validator

import (
	"strings"
	"testing"
)

func checkValidation(t *testing.T, gotErr error, expectErr bool) {
	t.Helper()

	if (gotErr != nil) != expectErr {
		t.Fatalf("expected error=%v got=%v", expectErr, gotErr)
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "too short",
			input:     "m",
			expectErr: true,
		},
		{
			name:      "minimum valid length",
			input:     "mo",
			expectErr: false,
		},
		{
			name:      "maximum valid length",
			input:     strings.Repeat("a", 50),
			expectErr: false,
		},
		{
			name:      "above maximum length",
			input:     strings.Repeat("a", 51),
			expectErr: true,
		},
		{
			name:      "normal name",
			input:     "Mohamed Elmorsy",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)

			checkValidation(t, err, tt.expectErr)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "missing at symbol",
			input:     "mkdfaddff.com",
			expectErr: true,
		},

		{
			name:      "invalid format",
			input:     "mldl",
			expectErr: true,
		},
		{
			name:      "short valid email",
			input:     "m@gmail.com",
			expectErr: false,
		},
		{
			name:      "normal valid email",
			input:     "mohamed@gmail.com",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)

			checkValidation(t, err, tt.expectErr)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "too short",
			input:     "2232j",
			expectErr: true,
		},
		{
			name:      "missing uppercase letter",
			input:     "mohamed3m4",
			expectErr: true,
		},
		{
			name:      "missing lowercase letter",
			input:     "MEIO@#MK4",
			expectErr: true,
		},
		{
			name:      "missing digit",
			input:     "MOhamedmkdfj$",
			expectErr: true,
		},
		{
			name:      "missing special character",
			input:     "Mohamedelmorsy3343",
			expectErr: true,
		},
		{
			name:      "valid password",
			input:     "Mohamed#@Elmorsy553",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)

			checkValidation(t, err, tt.expectErr)
		})
	}
}
