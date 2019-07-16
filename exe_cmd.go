package main

import (
	"fmt"

	"os/exec"
	"strings"
)

func exe_cmd(cmd string) ([]byte, error) {
	fmt.Println("cmd: ", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		return out, err
	}
	return out, err
}
