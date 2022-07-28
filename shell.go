package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

var prevCWD string
var cwd string

var userName string
var homeDir string

// host = strings.ReplaceAll(host, ".local", "")

func main() {
	fmt.Fprint(os.Stdout, "\\033]0;Gosh Shell\\007\r")
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	getUser()
	setCWD(&homeDir)
	reader := bufio.NewReader(os.Stdin)
	for {
		cwd, _ = os.Getwd()
		fmt.Printf("%s@%s %s $ ", userName, host, cwd)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err = execInput(&input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func getUser() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userName = usr.Username
	homeDir = usr.HomeDir
}

func setCWD(path *string) error {
	if err := os.Chdir(*path); err != nil {
		return fmt.Errorf("%s", err)
	}
	prevCWD = cwd
	cwd = *path
	return nil
}

func setChmod(path *string, mode *uint32) error {
	if err := os.Chmod(*path, os.FileMode(*mode)); err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

func execInput(input *string) error {
	*input = strings.TrimSuffix(*input, "\n")
	args := strings.Split(*input, " ")
	switch args[0] {
	case "cd":
		err := setCWD(&args[1])
		if err != nil {
			return err
		}
		return nil
	case "chmod":
		var mode uint32
		u64, _ := strconv.Atoi(args[1])
		mode = uint32(u64)
		err := setChmod(&args[0], &mode)
		if err != nil {
			return err
		}
		return nil
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
}
