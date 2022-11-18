package sftp

import (
	"bufio"
	"testing"
)

func TestRead(t *testing.T) {
	remoteFile := RemoteFile{
		Host: "127.0.0.1",
		Port: 10022,
		Url:  "/data/make/Makefile",
		User: "root",
	}
	read, err := Read(remoteFile)
	if err != nil {
		t.Logf("read error %v", err)
		return
	}
	reader := bufio.NewReader(read)
	line, err := reader.ReadString('\r')
	t.Logf("file : %s", line)
	read.Close()
}
