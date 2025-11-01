package phplint

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLintTestData(t *testing.T) {
	if os.Getenv("NIX_CC") != "" {
		t.Skip("Downloading does not work in Nix build")
	}

	supportedPHPVersions := []string{"8.1", "8.2", "8.3"}

	for _, version := range supportedPHPVersions {
		errors, err := LintFolder(t.Context(), version, "testdata")

		assert.NoError(t, err)

		assert.Len(t, errors, 1)

		assert.Equal(t, "invalid.php", errors[0].File)
	}
}
