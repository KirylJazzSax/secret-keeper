package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVaLidateEmail(t *testing.T) {
	cases := make(map[string]bool)
	cases["notvalid"] = false
	cases["valid@mail.com"] = true

	for email, isValid := range cases {
		if isValid {
			require.NoError(t, ValidateEmail(email))
			continue
		}

		require.Error(t, ValidateEmail(email))
	}
}
