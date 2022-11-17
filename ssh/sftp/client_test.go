package sftp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConnectKnownHost - connect to linux server in docker.
// command :
// docker run --privileged=true -itd --name ubuntu-2004-golang-ssh -p 18070:8070 -p 17060:7060 -p 10022:22 -v /D/coding/algorithm/distributed_system/deployment/single_server/data:/data ubuntu:20.04-golang /sbin/init
func TestConnectKnownHost(t *testing.T) {
	dir, _ := os.UserHomeDir()
	os.Setenv("HOME", dir)
	t.Logf("id_rsa : %s \n", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	t.Logf("known_hosts : %s \n", filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	client, err := ConnectKnownHost(
		"127.0.0.1",
		10022,
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
		filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"),
		"root")
	if err != nil {
		t.Logf("connect error : %v", err)
		return
	}
	root, err := client.Glob("/data/projects/**/*.go")
	if err != nil {
		t.Logf("getwd error : %v", err)
		return
	}
	t.Logf("root : %s \n", strings.Join(root, " || "))
	t.Logf("client : %s \n", client.String())
	client.Close()
}
