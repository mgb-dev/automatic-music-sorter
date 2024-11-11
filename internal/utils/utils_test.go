package utils

import (
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	// all lower case
	t.Run("already_lowercase", func(t *testing.T) {
		newDirName := "the greatest to ever exist"
		expected := strings.ReplaceAll(newDirName, " ", "-")
		obtained := NormalizeDirName(newDirName)

		if expected != obtained {
			t.Errorf("expected %v, obtained %v", expected, obtained)
		}
	})
	// no spaces
	t.Run("no_spaces_to_replace", func(t *testing.T) {
		newDirName := "SomETHING_REALLY_odd"
		expected := "something_really_odd"
		obtained := NormalizeDirName(newDirName)

		if expected != obtained {
			t.Errorf("expected %v, obtained %v", expected, obtained)
		}
	})
	// accented letters
	t.Run("letters_with_accents", func(t *testing.T) {
		newDirName := "SÓASDASÄSDÄSDÁSD Ë"
		expected := "sóasdasäsdäsdásd-ë"
		obtained := NormalizeDirName(newDirName)

		if expected != obtained {
			t.Errorf("expected %v, obtained %v", expected, obtained)
		}
	})
}
