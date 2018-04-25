package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var r = color.HiRedString
var y = color.HiYellowString
var g = color.HiGreenString

var em = color.New(color.ReverseVideo, color.Bold).Sprintf

var warnings int

func main() {
	em("checking for spew calls...")
	content := execute("git diff | cat")
	if strings.Contains(content, "spew.Dump") {
		fmt.Printf("%s\n", r("✘ found spew.Dump calls in your code. Aborting."))
		os.Exit(1)
	}

	ask("Is your %v up to date?", em(g("README.md")))
	ask("Is your %v up to date?", g("swagger.yaml"))
	ask("Is your %s up to date?", g("Postman collection"))
	ask("Is your %s shared and accessible?", g("relevant documentation"))
	checkWarnings()
}

func checkWarnings() {
	if warnings > 2 {
		fmt.Printf("---\n%s\n", r("✘ %d warnings. You are not ready for a merge request.", warnings))
		os.Exit(1)
	}
}

func ask(format string, args ...interface{}) {
	const prefix = "- "
	const suffix = " | (y/N)\n"

	format = fmt.Sprintf("%s %s %s", prefix, format, suffix)

	fmt.Printf(format, args)
	var answer string
	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	// reader := bufio.NewReader(tty)
	_, err = fmt.Fscanln(tty, &answer)
	// _, err := fmt.Scanln(&answer)
	if err != nil {
		fmt.Println(r("invalid answer..."))
		os.Exit(1)
	}
	if !strings.EqualFold(answer, "y") {
		fmt.Printf("\033[1A")
		fmt.Printf("\r%s %s\n", r("✘"), y("WARNING: You may not be ready for a merge request."))
		warnings++

		return
	}
	fmt.Printf("\033[1A")        // <- This one jumps back 1 line
	fmt.Printf("\r%s\n", g("✓")) // <- this positions the cursor at the begining of the line
}

// vex the command
func execute(command string) string {
	// invoke the naked shell
	cmd := exec.Command("/bin/sh", "-c", command)
	// only pipe the outputs
	var b []byte
	buf := bytes.NewBuffer(b)
	cmd.Stdin = os.Stdin
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	// run and log any errors
	if err := cmd.Run(); err != nil {
		log.Fatalf("cannot execute command: %v", err)
	}

	return buf.String()
}
