package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// newClient creates a new sftp/ssh client
func newSSH(ip, username, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// executeSSH executes commands on a remote server.
func executeSSH(ip, username, password string, commands []string) error {
	client, err := newSSH(ip, username, password)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// get input stream + bind outputs
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	err = session.Shell()
	if err != nil {
		return err
	}

	// send the commands
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return err
		}
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}
