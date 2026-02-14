// Various function for dealing with recipes.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// Try to unindent a recipe, so that it begins an column 0. (This is mainly for
// recipes in python, or other indentation-significant languages.)
func stripIndentation(s string, mincol int) string {
	// trim leading whitespace
	reader := bufio.NewReader(strings.NewReader(s))
	output := ""
	for {
		line, err := reader.ReadString('\n')
		col := 0
		i := 0
		for i < len(line) && col < mincol {
			c, w := utf8.DecodeRuneInString(line[i:])
			if strings.IndexRune(" \t\n", c) >= 0 {
				col += 1
				i += w
			} else {
				break
			}
		}
		output += line[i:]

		if err != nil {
			break
		}
	}

	return output
}

// Indent each line of a recipe.
func printIndented(out io.Writer, s string, ind int) {
	indentation := strings.Repeat(" ", ind)
	reader := bufio.NewReader(strings.NewReader(s))
	firstline := true
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if !firstline {
				io.WriteString(out, indentation)
			}
			io.WriteString(out, line)
		}
		if err != nil {
			break
		}
		firstline = false
	}
}

// Execute a recipe.
func dorecipe(target string, u *node, e *edge, rs *ruleSet, dryrun bool) (bool, int, string) {
	vars := make(map[string][]string)

	// Copy all global variables like $MKSHELL and environment variables
	for k, v := range rs.vars {
		vars[k] = v
	}

	vars["target"] = []string{target}
	if e.r.ismeta {
		if e.r.attributes.regex {
			for i := range e.matches {
				vars[fmt.Sprintf("stem%d", i)] = e.matches[i : i+1]
			}
		} else {
			vars["stem"] = []string{e.stem}
		}
	}

	// TODO: other variables to set
	// alltargets
	// newprereq

	vars["pid"] = []string{fmt.Sprintf("%d", os.Getpid())}

	prereqs := make([]string, 0)
	for i := range u.prereqs {
		if u.prereqs[i].r == e.r && u.prereqs[i].v != nil {
			prereqs = append(prereqs, u.prereqs[i].v.name)
		}
	}
	vars["prereq"] = prereqs

	input := expandRecipeSigils(e.r.recipe, vars)
	sh := "sh"
	args := []string{"-e"}

	if mkshell, ok := rs.vars["MKSHELL"]; ok && len(mkshell) > 0 {
		sh = mkshell[0]
		args = mkshell[1:]
	}

	if len(e.r.shell) > 0 {
		sh = e.r.shell[0]
		args = e.r.shell[1:]
	}

	mkPrintRecipe(target, input, e.r.attributes.quiet)

	// Export variables to the child shell environment exactly like Plan 9 mk
	env := os.Environ()
	for k, v := range vars {
		env = append(env, fmt.Sprintf("%s=%s", k, strings.Join(v, " ")))
	}

	// Explicitly inject MKSHELL=sh so the shell can evaluate $MKSHELL natively
	env = append(env, fmt.Sprintf("MKSHELL=%s", sh))

	if dryrun {
		return true, 0, input
	}

	_, success, exitcode := subprocessExit(
		sh,
		args,
		input,
		false,
		env)

	return success, exitcode, input
}

// subprocess executes the named program with args, feeding input to stdin. If
// captureOut is true, it captures and returns the program's stdout. Returns
// (output, success) where success is true if the process exits with code 0.
func subprocess(program string,
	args []string,
	input string,
	captureOut bool,
	env []string) (string, bool) {
	out, succ, _ := subprocessExit(program, args, input, captureOut, env)
	return out, succ
}

// subprocessExit executes the named program with args, feeding input to its
// stdin. If capture_out is true the program's stdout is captured and returned.
// It returns (stdout, success, exitCode) where success is true when
// exitCode == 0
func subprocessExit(program string,
	args []string,
	input string,
	capture_out bool,
	env []string) (string, bool, int) {
	program_path, err := exec.LookPath(program)
	if err != nil {
		log.Fatal(err)
	}

	proc_args := []string{program}
	proc_args = append(proc_args, args...)

	stdin_pipe_read, stdin_pipe_write, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	attr := os.ProcAttr{Files: []*os.File{stdin_pipe_read, os.Stdout, os.Stderr}, Env: env}

	output := make([]byte, 0)
	capture_done := make(chan bool)
	if capture_out {
		stdout_pipe_read, stdout_pipe_write, err := os.Pipe()
		if err != nil {
			log.Fatal(err)
		}

		attr.Files[1] = stdout_pipe_write

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdout_pipe_read.Read(buf)

				if err == io.EOF && n == 0 {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				output = append(output, buf[:n]...)
			}

			capture_done <- true
		}()
	}

	proc, err := os.StartProcess(program_path, proc_args, &attr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		_, err := stdin_pipe_write.WriteString(input)
		if err != nil {
			log.Fatal(err)
		}

		err = stdin_pipe_write.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	state, err := proc.Wait()

	if attr.Files[1] != os.Stdout {
		attr.Files[1].Close()
	}

	if err != nil {
		log.Fatal(err)
	}

	// wait until stdout copying in finished
	if capture_out {
		<-capture_done
	}

	return string(output), state.Success(), state.ExitCode()
}
