package main

import (
	"io"
	"os/exec"
)

type GoSubprocess struct {
	cmd *exec.Cmd
	r   *io.PipeReader
	w   *io.PipeWriter
}

func executeGoSubprocess(id string) (*GoSubprocess, error) {
	args := []string{
		"run", "./sub/process.go", "-id", id,
	}

	r, w := io.Pipe()

	cmd := exec.Command("go", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &GoSubprocess{cmd: cmd, r: r, w: w}, nil
}
