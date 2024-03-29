package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type scriptResult struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	ExitCode int    `json:"exit_code"`
}

var url string

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	var cmd exec.Cmd

	switch len(os.Args) {
	case 1, 2:
		panic("USAGE: ./runner name executable [args]")
	default:
		cmd = *exec.Command("/usr/bin/caffeinate", os.Args[2:]...)
	}

	name := os.Args[1]

	var outbuf, errbuf bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	var (
		stdout, stderr string
	)
	stdout = outbuf.String()
	stderr = errbuf.String()

	fmt.Fprint(os.Stdout, stdout)
	fmt.Fprint(os.Stderr, stderr)

	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exitCode = e.ExitCode()
		} else {
			panic(err)
		}
	}

	b, err := json.Marshal(scriptResult{
		Name:     name,
		Message:  stderr,
		ExitCode: exitCode,
	})
	if err != nil {
		panic(err)
	}

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*10)
	defer cncl()

	resp, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
