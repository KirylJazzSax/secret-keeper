package validation

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/require"
)

func TestVaLidateEmail(t *testing.T) {
	v := do.MustInvoke[Validator](nil)
	cases := make(map[string]bool)
	cases["notvalid"] = false
	cases["valid@mail.com"] = true

	for email, isValid := range cases {
		if isValid {
			require.NoError(t, v.ValidateEmail(email))
			continue
		}

		require.Error(t, v.ValidateEmail(email))
	}
}
