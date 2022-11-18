package sftp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var lock sync.RWMutex

// clientCache - key: ip of remote server.
var clientCache map[string]*Client

func init() {
	clientCache = make(map[string]*Client)
}

const IpRegex = `^([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$`

func IsIP(address string) bool {
	compile := regexp.MustCompile(IpRegex)
	return compile.MatchString(address)
}

// RemoteFile - used to read/write remote file by sftp/ssh
type RemoteFile struct {
	Host string
	Port int
	Url  string
	User string
}

func (file RemoteFile) String() string {
	return fmt.Sprintf("%s@%s:%d[%s]", file.User, file.Host, file.Port, file.Url)
}
func (file RemoteFile) IsValid() bool {
	if file.Port == 0 {
		file.Port = 22
	}
	file.Host = strings.TrimSpace(file.Host)
	file.Url = strings.TrimSpace(file.Url)
	file.User = strings.TrimSpace(file.User)
	return len(file.Host) != 0 &&
		len(file.Url) != 0 &&
		len(file.User) != 0
}

// Read - read remote file by sftp with paired keys.
func Read(file RemoteFile) (io.ReadCloser, error) {
	if !file.IsValid() {
		return nil, errors.New(fmt.Sprintf("file[%s] is invalid", file))
	}
	if !IsIP(file.Host) {
		// 此时可能是主机名
		addrs, err := net.LookupHost(file.Host)
		if err != nil {
			return nil, fmt.Errorf("lookup for host[%s] error[%w]", file.Host, err)
		}
		file.Host = addrs[0]
	}
	client := getClient(file)
	if client == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("can't get home directory, %w", err)
		}
		privateKeyFile := filepath.Join(home, ".ssh", "id_rsa")
		knownhostsFile := filepath.Join(home, ".ssh", "known_hosts")
		client, err = ConnectKnownHost(file.Host, file.Port, privateKeyFile, knownhostsFile, file.User)
		if err != nil {
			return nil, fmt.Errorf("connect to host[%s] error[%w]", file.Host, err)
		}
		putClient(client)
	}

	reader, err := client.Open(file.Url)
	if err != nil {
		return nil, fmt.Errorf("read remote file[%s] error[%w]", file, err)
	}
	return reader, nil
}

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

func Write(file RemoteFile) (ReadWriteSeekCloser, error) {
	if !file.IsValid() {
		return nil, errors.New(fmt.Sprintf("file[%s] is invalid", file))
	}
	if !IsIP(file.Host) {
		// 此时可能是主机名
		addrs, err := net.LookupHost(file.Host)
		if err != nil {
			return nil, fmt.Errorf("lookup for host[%s] error[%w]", file.Host, err)
		}
		file.Host = addrs[0]
	}
	client := getClient(file)
	if client == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("can't get home directory, %w", err)
		}
		privateKeyFile := filepath.Join(home, ".ssh", "id_rsa")
		knownhostsFile := filepath.Join(home, ".ssh", "known_hosts")
		client, err = ConnectKnownHost(file.Host, file.Port, privateKeyFile, knownhostsFile, file.User)
		if err != nil {
			return nil, fmt.Errorf("connect to host[%s] error[%w]", file.Host, err)
		}
		putClient(client)
	}
	openFile, err := client.OpenFile(file.Url, os.O_RDWR|os.O_APPEND)
	if err != nil {
		return nil, fmt.Errorf("open remote file[%s] error[%w]", file, err)
	}
	return openFile, nil
}

func putClient(client *Client) {
	lock.Lock()
	defer lock.Unlock()
	clientCache[client.Host()] = client
}
func getClient(file RemoteFile) *Client {
	lock.RLock()
	defer lock.RUnlock()
	if client, ok := clientCache[file.Host]; ok {
		return client
	}
	return nil
}
