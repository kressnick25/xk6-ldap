package ldap

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExamplesInputOutput runs same k6's scripts that we have in example folder
// it check that output contains/not contains cetane things
// it's not a real test, but it's a good way to check that examples are working
// between changes
//
// We also do use a convention that successful output should contain `level=info` (at least one info message from console.log), e.g.:
// INFO[0000] deciphered text == original text:  true       source=console
// and should not contain `level=error` or "Uncaught", e.g. outputs like:
// ERRO[0000] Uncaught (in promise) OperationError: length is too large  executor=per-vu-iterations scenario=default
func TestExamplesInputOutput(t *testing.T) {
	t.Parallel()

	outputShouldContain := []string{
		"output: -",
		"default: 1 iterations for each of 1 VUs",
		"1 complete and 0 interrupted iterations",
		"level=info", // at least one info message
	}

	outputShouldNotContain := []string{
		"Uncaught",
		"level=error", // no error messages
	}

	// List of the directories containing the examples
	// that we should run and check that they produce the expected output
	// and not the unexpected one
	// it could be a file (ending with .js) or a directory
	examples := []string{
		"./examples/example.js",
		"./examples/example-tls.js",
	}

	for _, path := range examples {
		list := getFiles(t, path)

		for _, file := range list {
			name := filepath.Base(file)
			file := file

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				cmd := exec.Command("./k6", "run", "-v", "--log-output=stdout", file) /* #nosec G204 */
				var stdout, stderr strings.Builder
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					panic(err)
				}

				for _, s := range outputShouldContain {
					assert.Contains(t, stdout.String(), s)
				}
				for _, s := range outputShouldNotContain {
					assert.NotContains(t, stdout.String(), s)
				}

				assert.Empty(t, stderr.String())
			})
		}
	}
}

func getFiles(t *testing.T, path string) []string {
	t.Helper()

	result := []string{}

	// If the path is a file, return it as is
	if strings.HasSuffix(path, ".js") {
		return append(result, path)
	}

	// If the path is a directory, return all the files in it
	list, err := os.ReadDir(path) //nolint:forbidigo // we read a directory
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	for _, file := range list {
		if file.IsDir() {
			continue
		}

		result = append(result, filepath.Join(path, file.Name()))
	}

	return result
}
