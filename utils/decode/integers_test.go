package decode

import (
	"log"
	"testing"
)

func TestToInteger(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	test1 := "i123e"

	result, err := ToInteger(test1)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	log.Printf("result: %d", result)

	if result != 123 {
		t.Errorf("Expected: 123, got: %d", result)
	}
}
