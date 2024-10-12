package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExamples(t *testing.T) {
	dir := "./examples"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("directory does not exist: %s", dir)
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() != "examples" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			t.Run(fmt.Sprintf("test_example_%s", info.Name()), func(t *testing.T) {
				runny := Runny{
					Config: Config{
						File:    path,
						Testing: true,
					},
				}
				runny.Run()
			})
		}

		return nil
	})

	if err != nil {
		t.Fatalf("error walking the directory: %v", err)
	}
}
