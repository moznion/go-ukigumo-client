package ukigumo

import (
	"fmt"
	"os"
	"os/exec"
)

func executeCommand(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}
	return string(out), nil
}

func executeCommandWithOutput(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}
	fmt.Printf("%s\n", out)
	return string(out), nil
}
