package main

import (
	"testing"
)

func TestMainStep1(t *testing.T) {
	t.Run("parseJSON2", func(t *testing.T) {
		parseJSON2("./tests/step1/valid.json")
		parseJSON2("./tests/step1/invalid.json")
	})
}

func TestMainStep2Valid(t *testing.T) {
	t.Run("parseJSON2", func(t *testing.T) {
		parseJSON2("./tests/step2/valid.json")
		parseJSON2("./tests/step2/valid2.json")
	})
}

func TestMainStep2Invalid(t *testing.T) {
	t.Run("parseJSON2", func(t *testing.T) {
		expected1 := "unexpected character } on line 1"
		expected2 := "unexpected character k on line 3"
		_, err := parseJSON2("./tests/step2/invalid.json")
		if err.Error() != expected1 {
			t.Errorf("Expected error to be %s, got %v", expected1, err)
		}
		_, err2 := parseJSON2("./tests/step2/invalid2.json")
		if err2.Error() != expected2 {
			t.Errorf("Expected error to be %s, got %v", expected2, err2)
		}
	})
}

func TestMainStep3Valid(t *testing.T) {
	t.Run("parseJSON2", func(t *testing.T) {
		parseJSON2("./tests/step3/valid.json")
	})
}

func TestMainStep3Invalid(t *testing.T) {
	t.Run("parseJSON2", func(t *testing.T) {
		expected := "unexpected character F on line 3"
		_, err := parseJSON2("./tests/step3/invalid.json")
		if err.Error() != expected {
			t.Errorf("Expected error to be %s, got %v", expected, err)
		}
	})
}
