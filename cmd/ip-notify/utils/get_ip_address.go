package utils

import (
	"os/exec"
	"strings"
)

func GetIPAddress() ([]string, error) {
	out, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		return nil, err
	}
	ips := strings.Fields(string(out))
	return ips, nil
}
