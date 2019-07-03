package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// parseCams parses camString to a map of cameras w/ ips
// TODO: check for invalid string splits
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

// parseCommands parses commands from command line args
// TODO: validate commands against argument list
func parseCommands(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("No commands detected! %v", args)
	}
	var command string
	for i, x := range args {
		if i > 1 {
			command = command + " " + x
		}
	}
	return command, nil
}

// executeCommands executes parsed commands; currently only supports single commands
func executeCommands(command string, camMap map[string]string, selectedCam, username, password string) error {

	words := strings.Fields(command)
	customCommand := isCustomCommand(words[0])
	mappedCommand := isMappedCommand(words[0])

	if customCommand {
		// exeute custom command
		fmt.Printf("Executing a custom command on %s ... \n", selectedCam)
		err := executeCustomCommand(selectedCam, username, password, words[0])
		if err != nil {
			return err
		}
		fmt.Printf("Command: %s ... Success!\n", words[0])
	} else if mappedCommand {
		// execute mapped command
		fmt.Printf("Executing a mapped command on %s ... \n", selectedCam)
		cmds := []string{strings.Replace(command, words[0], commandMap[words[0]], 1), "exit"}
		err := executeSSH(selectedCam, username, password, commands)
		if err != nil {
			return err
		}
		fmt.Printf("Command: %s ... Success!\n", cmds[0])
	} else {
		// execute literal command
		fmt.Printf("Executing a literal command on %s ... \n", selectedCam)
		cmds := []string{command, "exit"}
		err := executeSSH(selectedCam, username, password, commands)
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

// isMappedCommand cheks if command is in commandMap
func isMappedCommand(command string) bool {
	for k := range commandMap {
		if k == command {
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

	// parse commands
	command, err := parseCommands(args)
	if err != nil {
		log.Fatalln(err)
	}

	// execute commands
	err = executeCommands(command, cams, selectedCam, username, password)
	if err != nil {
		log.Fatalln(err)
	}
}
