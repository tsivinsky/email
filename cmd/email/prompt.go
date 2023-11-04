package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func prompt(qs string) (string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print(qs)
	ans, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(ans), nil
}

func mustPrompt(qs string) string {
	ans, err := prompt(qs)
	if err != nil {
		panic(err)
	}

	return ans
}

func getEmailBody() ([]byte, error) {
	f, err := os.CreateTemp(os.TempDir(), "email-body")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cmd := exec.Command("nvim", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	os.Remove(f.Name())

	return data, nil
}
