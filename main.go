package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ExecuteSSH executes commands on a remote server.
func ExecuteSSH(ip, username, password string, commands []string) {
	// create ssh config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Fatal(err)
	}
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// get input stream + bind outputs
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	// send the commands
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// ExecuteHTTP executes commands on a fang via HTTP.
// func ExecuteHTTP(ip, command string) {
// 	var httpAddress string
// 	if len(ip) > 0 && len(command) > 0 {
// 		httpAddress = fmt.Sprintf("http://%s/cgi-bin/action.cgi?cmd=%s", ip, command)
// 	}
// 	resp, err := http.Get(httpAddress)
// 	if err != nil {
// 		log.Fatalf("Failed to execute http command (%s, %s): %v", ip, command, err)
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to read http response body(%s): %v", httpAddress, err)
// 	}
// 	fmt.Println(httpAddress)
// 	fmt.Println(string(body))
// }

// parseCams parses camString to a map of cameras w/ ips
func parseCams(camString string) (map[string]string, error) {
	var camMap = make(map[string]string)
	cams := strings.Split(camString, "//")
	for _, cam := range cams {
		temp := strings.Split(cam, "/")
		camMap[temp[0]] = temp[1]
	}
	return camMap, nil
}

func main() {

	// setup
	_, userBool := os.LookupEnv("GFANG_USER")
	if !userBool {
		log.Fatalln("GFANG_USER environment variable not found!")
	}
	_, passBool := os.LookupEnv("GFANG_PASS")
	if !passBool {
		log.Fatalln("GFANG_PASS environment variable not found!")
	}
	camString, camsBool := os.LookupEnv("GFANG_CAMS")
	if !camsBool {
		log.Fatalln("GFANG_CAMS environment variable not found!")
	}
	cams := parseCams(camString)

	args := os.Args
	fmt.Println(len(args))

	// commands := []string{
	// 	"motor up 100",
	// 	//"night_mode status",
	// 	"exit",
	// }

	//ExecuteHTTP(fangs[0], "motor_right")

	//ExecuteSSH(fangs[1], username, password, commands)

	// for _, fang := range fangs {
	// 	ExecuteSSH(fang, username, password, commands)
	// }

}
