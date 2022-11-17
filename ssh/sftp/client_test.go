package sftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConnectKnownHost(t *testing.T) {
	log.Printf("id_rsa : %s \n", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	log.Printf("known_hosts : %s \n", filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	client, err := ConnectKnownHost(
		"127.0.0.1",
		10022,
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
		filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"),
		"root")
	if err != nil {
		fmt.Printf("connect error : %v", err)
		return
	}
	root, err := client.Glob("/data/projects/**/*.go")
	if err != nil {
		fmt.Printf("getwd error : %v", err)
		return
	}
	fmt.Printf("root : %s \n", strings.Join(root, " || "))
	fmt.Printf("client : %s \n", client.String())
	client.Close()
}
