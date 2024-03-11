package main

import (
	"os"
	"testing"
)

// TODO
func TestCopyExampleApp(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = initGomxApp("test")
	if err != nil {
		t.Fatal(err)
	}
}
