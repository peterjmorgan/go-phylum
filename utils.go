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

	start := recPat.FindStringIndex(desc)//[0]
	if len(start) > 0 {
		offset := start[0]
		result = desc[offset:]
	}
	//result = desc[start:]

	return result, nil
}
