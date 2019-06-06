package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"golang.org/x/crypto/ssh"
)

// getSnapshot retrieves a snapshot from a cam
func getSnapshot(username, password, ip string) error {

	// connect to cam
	var conn *ftp.ServerConn
	conn, err := ftp.Dial(ip+":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		// enable ftp and retry
		err = enableFTP(username, password, ip)
		if err != nil {
			log.Fatal(err)
		}
		conn, err = ftp.Dial(ip+":21", ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = conn.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	// retrieve snapshot
	resp, err := conn.Retr("/tmp/snapshot.jpg")
	if err != nil {
		log.Fatalf("snapshot not found! %v", err)
	}
	defer resp.Close()
	oFile, err := os.Create("snapshot.jpg")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	for {
		n, err := resp.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if _, err := oFile.Write(buf[:n]); err != nil {
			log.Fatal(err)
		}
	}
	err = conn.Quit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// enableFTP enables FTP on cam
func enableFTP(username, password, ip string) error {
	commands := []string{
		"bftpd -D &",
		"exit",
	}
	err := executeSSH(ip, username, password, commands)
	if err != nil {
		return err
	}
	return nil
}

// newClient creates a new sftp/ssh client
func newClient(ip, username, password string) (*ssh.Client, error) {
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
	client, err := newClient(ip, username, password)
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
