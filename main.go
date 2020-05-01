package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// parseCams parses camString to a map of cameras w/ ips
// TODO: check for invalid string splits using regex
func parseCams(camString string) (map[string]string, error) {
	var camMap = make(map[string]string)
	cams := strings.Split(camString, "//")
	for _, cam := range cams {
		temp := strings.Split(cam, "/")
		camMap[temp[1]] = temp[0]
	}
	return camMap, nil
}

// findCam tries to find cam in camMap
// is cam a number?
func findCam(camMap map[string]string, cam string) (string, error) {
	var selectedCam string
	if camNumber, err := strconv.Atoi(cam); err == nil {
		// it is a number
		count := 0
		for x := range camMap {
			if camNumber == count {
				selectedCam = x
			}
		}
	} else {
		// it is a string
		for x := range camMap {
			if x == cam {
				selectedCam = x
			}
		}
	}
	if selectedCam == "" {
		return "", fmt.Errorf("Failed to find cam! %s", cam)
	}
	return camMap[selectedCam], nil
}

// executeCommands executes parsed commands; currently only supports single commands
func executeCommands(commands []string, camMap map[string]string, selectedCam, username, password string) error {

	customCommand := isCustomCommand(commands[1])

	if customCommand {
		// execute custom command
		fmt.Printf("Executing a custom command on %s ... \n", selectedCam)
		err := executeCustomCommand(selectedCam, username, password, commands[1])
		if err != nil {
			return err
		}
		fmt.Printf("Command: %s ... Success!\n", commands[1])
	} else {
		// execute literal command
		fmt.Printf("Executing a literal command on %s ... \n", selectedCam)
		var command string
		for i, x := range commands {
			if i > 1 {
				command = command + " " + x
			}
		}
		cmds := []string{command, "exit"}
		err := executeSSH(selectedCam, username, password, cmds)
		if err != nil {
			return err
		}
		fmt.Printf("Command: %s ... Success!\n", cmds[0])
	}

	return nil
}

// isCustomCommand checks if command is in customCommandsMap
func isCustomCommand(command string) bool {
	for _, cmd := range customCommands {
		if cmd == command {
			return true
		}
	}
	return false
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("Insufficient command line arguments")
		// TODO: print example usage string
		// TODO: enable command selection
	}
	args := os.Args

	// fetch authentication
	username, userBool := os.LookupEnv("GFANG_USER")
	if !userBool {
		log.Fatalln("GFANG_USER environment variable not found!")
	}
	password, passBool := os.LookupEnv("GFANG_PASS")
	if !passBool {
		log.Fatalln("GFANG_PASS environment variable not found!")
	}

	// fetch,parse, and find cam(s)
	camString, camsBool := os.LookupEnv("GFANG_CAMS")
	if !camsBool {
		log.Fatalln("GFANG_CAMS environment variable not found!")
	}
	cams, err := parseCams(camString)
	if err != nil {
		log.Fatalf("Failed to parse cam list! %s", camString)
	}
	selectedCam, err := findCam(cams, args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// execute commands
	err = executeCommands(args, cams, selectedCam, username, password)
	if err != nil {
		log.Fatalln(err)
	}
}
