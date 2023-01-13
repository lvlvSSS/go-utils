package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	host := "QHNKgQPf5lRaM"
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey for %s", host)
	}
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		fmt.Println(err)
		return
	}

	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(singer)},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:22", config)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer session.Close()

	var exitcode int

	output, err := session.CombinedOutput(`echo $PATH`)
	if err != nil {

		if ins, ok := err.(*ssh.ExitError); ok {
			exitcode = ins.ExitStatus()
		} else {
			exitcode = ins.ExitStatus()
		}
	}
	fmt.Println(string(output), exitcode)

}
