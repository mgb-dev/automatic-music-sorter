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
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// no spaces
	t.Run("no_spaces_to_replace", func(t *testing.T) {
		newDirName := "SomETHING_REALLY_odd"
		expected := "something_really_odd"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// accented letters
	t.Run("letters_with_accents", func(t *testing.T) {
		newDirName := "SÓASDASÄSDÄSDÁSD Ë"
		expected := "sóasdasäsdäsdásd-ë"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// remove colon
	t.Run("remove_colon", func(t *testing.T) {
		newDirName := "the name: is the primegean"
		expected := "the-name-is-the-primegean"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// semi-colon
	t.Run("remove_semi_colon", func(t *testing.T) {
		newDirName := "the name is; the primegean"
		expected := "the-name-is-the-primegean"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// comma
	t.Run("remove_comma", func(t *testing.T) {
		newDirName := "the name, is the primegean"
		expected := "the-name-is-the-primegean"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
	// fwd_slash
	t.Run("remove_forward_slash", func(t *testing.T) {
		newDirName := "the name is /the primegean"
		expected := "the-name-is-the-primegean"
		obtained, err := NormalizeDirName(newDirName)
		if err != nil {
			t.Errorf("Something went wrong: %s", err)
		}

		if expected != obtained {
			t.Errorf("expected: %v, obtained: %v\n", expected, obtained)
		}
	})
}
