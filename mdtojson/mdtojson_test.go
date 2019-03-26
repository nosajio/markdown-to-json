package mdtojson

import (
	"strings"
	"testing"
)

func TestProcessRepo(t *testing.T) {
	t.Run("ProcessRepo(<URL>, <DIR>)", func(t *testing.T) {
		testRepoURL := "https://github.com/nosajio/writing"
		testRepoDIR := "/tmp/posts"
		json, err := ProcessRepo(testRepoURL, testRepoDIR)

		if err != nil {
			t.Errorf("ProcessRepo(%s, %s) errored while processing: %s", testRepoURL,
				testRepoDIR, err.Error())
		}

		if json == "" {
			t.Errorf("ProcessRepo(%s. %s) returned an empty JSON string", testRepoURL,
				testRepoDIR)
		}

		if !strings.HasPrefix(json, "[{") {
			t.Errorf("ProcessRepo(%s, %s) didn't return a valid JSON array:\n %s",
				testRepoURL, testRepoDIR, json)
		}
	})
}
