package decode

import (
	"log"
	"testing"
)

func TestString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	test1 := "5:hello"

	result, err := ToString(test1)

	log.Printf("result: %s", result)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if result != "hello" {
		t.Errorf("Expected: hello, got: %s", result)
	}
}
