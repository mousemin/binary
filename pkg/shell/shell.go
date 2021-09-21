package shell

import "os/exec"

func Mv(args ...string) error {
	cmd := exec.Command("mv", args...)
	_, err := cmd.Output()
	return err
}