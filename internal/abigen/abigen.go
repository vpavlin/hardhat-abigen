package abigen

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Run(in string, typ string, pkg string, output string) error {
	args := []string{
		"-abi", in,
		"-out", output,
		"-type", typ,
		"-pkg", pkg,
	}

	cmd := exec.Command("abigen", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to generate bindings: %s, stderr: %s", err, stdErr.String())
	}

	return nil
}
