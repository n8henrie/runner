package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type scriptResult struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	ExitCode int    `json:"exit_code"`
}

var url string

func main() {
	var (
		exitCode int
		message  string
	)
	var cmd exec.Cmd

	switch len(os.Args) {
	case 1, 2:
		panic("USAGE: ./runner name executable [args]")
	case 3:
		cmd = *exec.Command(os.Args[2])
	default:
		cmd = *exec.Command(os.Args[2], os.Args[3:]...)
	}
	name := os.Args[1]
	var outbuf, errbuf bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()

	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exitCode = e.ExitCode()
			message = errbuf.String()
		} else {
			panic(err)
		}
	}

	b, err := json.Marshal(scriptResult{
		Name:     name,
		Message:  message,
		ExitCode: exitCode,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}
	os.Exit(exitCode)
}
