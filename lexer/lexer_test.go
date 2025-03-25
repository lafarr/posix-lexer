package lexer

import (
	"testing"
)

func equalSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, token := range a {
		if token != b[idx] {
			return false
		}
	}
	return true
}

func sliceAsString(a []string) string {
	tokensAsString := "["
	for _, token := range a {
		tokensAsString += token + ", "
	}
	return tokensAsString[:len(tokensAsString)-2] + "]"
}

func TestTokenization(t *testing.T) {

}
