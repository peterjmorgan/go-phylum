package phylum

import (
	"regexp"

	"github.com/pkg/errors"
)

func ExtractRemediation(issue IssuesListItem) (string, error) {
	var result string
	desc := issue.Description

	recPat := regexp.MustCompile(`### Recommendation`)

	if len(desc) < 1 {
		// no description
		return "", errors.New("No description found")
	}

	start := recPat.FindStringIndex(desc)[0]
	result = desc[start:]

	return result, nil
}
