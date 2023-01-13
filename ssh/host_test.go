package main

import (
	"bufio"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHost(t *testing.T) {
	host := "QHNKgQPf5lRaM+XO8vWCJgA4I6I"
	file, err := os.Open("./know_hosts")
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
}

func TestKnownHosts(t *testing.T) {
	known_hosts, _ := os.ReadFile("./know_hosts")
	marker, h, _, comment, rest, err := ssh.ParseKnownHosts(known_hosts)
	if err != nil {
		t.Errorf("error : %s", err)
		return
	}
	t.Logf("marker: %s \n", marker)
	t.Logf("hosts: %s \n", strings.Join(h, ","))
	t.Logf("comment: %s \n", comment)
	t.Logf("rest: %v \n", rest)
}

func TestSlice(t *testing.T) {
	methods := make([]ssh.AuthMethod, 4)
	t.Logf("%v \n", len(methods))
	t.Logf("%v \n", cap(methods))
}

func TestWaitGroup(t *testing.T) {
	out := &outgroup{}
	out.ingroup.wg.Add(1)
	go out.Done()
	out.Wait()
	t.Logf("Done")
	t.Logf("%p", &out.ingroup.wg)
}

type outgroup struct {
	ingroup
}

type ingroup struct {
	wg sync.WaitGroup
}

func (in *ingroup) Done() {
	defer in.wg.Done()
	time.Sleep(5 * time.Second)
}
func (in *ingroup) Wait() {
	in.wg.Wait()
}
