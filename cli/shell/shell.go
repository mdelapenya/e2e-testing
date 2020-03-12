package shell

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Execute executes a command in the machine the program is running
// - workspace: represents the location where to execute the command
// - command: represents the name of the binary to execute
// - args: represents the arguments to be passed to the command
func Execute(workspace string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args[0:]...)

	cmd.Dir = workspace

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"baseDir": workspace,
			"command": command,
			"args":    args,
			"error":   err,
			"stderr":  stderr.String(),
		}).Error("Error executing command")

		return "", err
	}

	return strings.Trim(out.String(), "\n"), nil
}

// Which checks if software is installed, else it aborts the execution
func Which(binary string) error {
	path, err := exec.LookPath(binary)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"binary": binary,
		}).Error("Required binary is not present")
		return err
	}

	log.WithFields(log.Fields{
		"binary": binary,
		"path":   path,
	}).Debug("Binary is present")
	return nil
}