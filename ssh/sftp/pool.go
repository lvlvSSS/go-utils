package sftp

import (
	"sync"
)

var lock sync.RWMutex
var clientCache map[string]*Client

func init() {
	clientCache = make(map[string]*Client)
}

// RemoteFile - used to read/write remote file by sftp/ssh
type RemoteFile struct {
	Host string
	Port int
	Url  string
	User string
}

/*func Read(file RemoteFile) (io.ReadCloser, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("can't get home directory, %w")err
	}
}*/
