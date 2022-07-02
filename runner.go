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

	var (
		exitCode int
		message  string
	)
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

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*10)
	defer cncl()

	resp, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}
}
