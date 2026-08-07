package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// mk is the path to the mk binary used during tests.
var mk = "mk"

// TestMain checks if the mk binary exists
func TestMain(m *testing.M) {
	if env := os.Getenv("MKBIN"); env != "" {
		mk = env
	} else if _, err := os.Stat("./mk"); err != nil {
		fmt.Fprintf(os.Stderr, "./mk :%v\n", err)
		os.Exit(1)
	} else {
		mk = "./mk"
	}

	path, err := exec.LookPath(mk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s :%v\n", mk, err)
		os.Exit(1)
	}
	mk = path

	os.Exit(m.Run())
}

// runMk runs the mk binary with the given mkfile, from the repo root, and
// returns stdout, stderr, and the exit code.
func runMk(t *testing.T, mkfile string) (stdout, stderr string, exitCode int) {
	t.Helper()

	// Always run from the repo root so that relative paths in mkfiles work.
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(mk, "-f", mkfile)
	cmd.Dir = repoRoot

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running mk: %v", err)
		}
	}

	return stdout, stderr, exitCode
}

// TestLines verifies that backslash-newline continuation propagates correctly
// through comment lines (Plan 9 mk compatibility).
func TestLines(t *testing.T) {
	stdout, stderr, exitCode := runMk(t, "test/lines.mk")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr: %s", exitCode, stderr)
	}
	// The LINES variable should join the two echo commands into one output line.
	if !strings.Contains(stdout, "line 1") || !strings.Contains(stdout, "line 2") {
		t.Errorf("expected output to contain 'line 1' and 'line 2', got: %q", stdout)
	}
}

// TestVarInclude verifies that variable expansion works in '<' include paths.
func TestVarInclude(t *testing.T) {
	stdout, stderr, exitCode := runMk(t, "test/var/include.mk")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "pid=") {
		t.Errorf("expected output to contain 'pid=', got: %q", stdout)
	}
}

// TestFailNow verifies that a recipe that exits non-zero causes mk to fail
// with exit code 1.
func TestFailNow(t *testing.T) {
	_, _, exitCode := runMk(t, "test/fail/now.mk")
	if exitCode != 1 {
		t.Fatal("expected exit code 1, got: " + strconv.Itoa(exitCode))
	}
}

// TestFailStatus verifies that mk shows the exit code returned by a recipe.
func TestFailStatus(t *testing.T) {
	_, stderr, exitCode := runMk(t, "test/fail/status.mk")
	if exitCode != 1 {
		t.Fatal("expected exit code 1, got: " + strconv.Itoa(exitCode))
	}
	if !strings.Contains(stderr, "status=exit(111)") {
		t.Errorf("expected output to contain 'status=exit(111)', got: %q",
			stderr)
	}
}

// TestFailE verifies that mk succeeds even if any command in a recipe fails
// with E attribute.
func TestFailE(t *testing.T) {
	_, stderr, exitCode := runMk(t, "test/fail/e.mk")
	if exitCode != 0 {
		t.Fatalf("expected exit 0 with E attribute, got %d\nstderr: %s",
			exitCode, stderr)
	}
}

// TestFailD verifies that mk deletes the target on failure with D attribute.
func TestFailD(t *testing.T) {
	// Clean up any leftover artifact from a previous run.
	os.Remove("test/fail/d")

	_, _, exitCode := runMk(t, "test/fail/d.mk")
	if exitCode != 1 {
		t.Fatal("expected exit 1 with D attribute, got: " +
			strconv.Itoa(exitCode))
	}
	// The target file should have been deleted.
	if _, err := os.Stat("test/fail/d"); !os.IsNotExist(err) {
		t.Error("expected 'test/fail/d' to be deleted after D failure")
	}
}
