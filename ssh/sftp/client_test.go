package sftp

import (
	"bytes"
	"net"
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

func TestHostHash(t *testing.T) {
	testHostHash(t, "172.17.0.3", "|1|v3atYVzYo9koRFMa3PMTdL/XK20=|9CU9l8fFS276gpt2EEmy7ToSoFY=")
}

func testHostHash(t *testing.T, hostname, encoded string) {
	typ, salt, hash, err := decodeHash(encoded)
	if err != nil {
		t.Fatalf("decodeHash: %v", err)
	}

	if got := encodeHash(typ, salt, hash); got != encoded {
		t.Errorf("got encoding %s want %s", got, encoded)
	}

	if typ != sha1HashType {
		t.Fatalf("got hash type %q, want %q", typ, sha1HashType)
	}

	got := hashHost(hostname, salt)
	if !bytes.Equal(got, hash) {
		t.Errorf("got hash %x want %x", got, hash)
	}
}

func TestHostIp(t *testing.T) {
	ip, _ := net.LookupIP("192.168.50.188")
	t.Logf("%v", ip)
	addrs, _ := net.LookupHost("192.168.50.188")
	t.Logf("%v", addrs)
}
