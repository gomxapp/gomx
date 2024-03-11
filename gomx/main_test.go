package main

import (
	"os"
	"testing"
)

func TestCopyExampleApp(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = copyExampleApp("foo_shop", dir)
	if err != nil {
		t.Fatal(err)
	}
}
