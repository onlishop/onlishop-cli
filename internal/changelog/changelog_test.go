package changelog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onlishop/onlishop-cli/internal/git"
)

func TestGenerateWithoutConfig(t *testing.T) {
	commits := []git.GitCommit{
		{
			Message: "feat: add new feature",
			Hash:    "1234567890",
		},
	}

	changelog, err := renderChangelog(commits, Config{
		VCSURL:   "https://github.com/FriendsOfOnlishop/FroshTools/commit",
		Template: defaultChangelogTpl,
	})

	assert.NoError(t, err)

	assert.Equal(t, "- [feat: add new feature](https://github.com/FriendsOfOnlishop/FroshTools/commit/1234567890)", changelog)
}

func TestTicketParsing(t *testing.T) {
	commits := []git.GitCommit{
		{
			Message: "NEXT-1234 - Fooo",
			Hash:    "1234567890",
		},
	}

	cfg := Config{
		Variables: map[string]string{
			"ticket": "^(NEXT-[0-9]+)",
		},
		Template: "{{range .Commits}}- [{{ .Message }}](https://issues.onlishop.com/issues/{{ .Variables.ticket }}){{end}}",
	}

	changelog, err := renderChangelog(commits, cfg)

	assert.NoError(t, err)
	assert.Equal(t, "- [NEXT-1234 - Fooo](https://issues.onlishop.com/issues/NEXT-1234)", changelog)
}

func TestIncludeFilters(t *testing.T) {
	commits := []git.GitCommit{
		{
			Message: "NEXT-1234 - Fooo",
			Hash:    "1234567890",
		},
		{
			Message: "merge foo",
			Hash:    "1234567890",
		},
	}

	cfg := Config{
		Pattern:  "^(NEXT-[0-9]+)",
		Template: defaultChangelogTpl,
	}

	changelog, err := renderChangelog(commits, cfg)

	assert.NoError(t, err)
	assert.Equal(t, "- [NEXT-1234 - Fooo](/1234567890)", changelog)
}
