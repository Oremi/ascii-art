package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runProgram(input string) (string, error) {
	c1 := exec.Command("go", "run", ".", input)
	c2 := exec.Command("cat", "-e")

	// Connect the pipe
	reader, err := c1.StdoutPipe()
	if err != nil {
		return "", err
	}
	c2.Stdin = reader

	// Start c1
	if err := c1.Start(); err != nil {
		return "", err
	}

	// CombinedOutput() runs c2, waits for it to finish, and collects output.
	// This will return once c1's stdout is closed (when c1 exits).
	out, err := c2.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	// Clean up c1
	if err := c1.Wait(); err != nil {
		return string(out), err
	}

	return string(out), nil
}

func TestAsciiArt(t *testing.T) {
	testDir := "tests"

	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("failed to read test directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") && !strings.Contains(file.Name(), "expected") {

			inputPath := filepath.Join(testDir, file.Name())
			expectedPath := filepath.Join(testDir, strings.Replace(file.Name(), ".txt", "_expected.txt", 1))

			inputBytes, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("failed reading input file: %v", err)
			}

			expectedBytes, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("failed reading expected file: %v", err)
			}

			input := strings.TrimSpace(string(inputBytes))
			expected := string(expectedBytes)

			t.Run(file.Name(), func(t *testing.T) {

				out, err := runProgram(input)
				if err != nil {
					t.Fatalf("program failed: %v", err)
				}

				if out != expected {
					t.Errorf("\nExpected:\n%s\n\nGot:\n%s", expected, out)
				}
			})
		}
	}

}
