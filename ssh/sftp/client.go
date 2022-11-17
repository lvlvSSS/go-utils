package sftp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client - this client could do some i/o for remote file by ssh/sftp
type Client struct {
	username       string
	certKeyFile    string
	knownHostsFile string
	password       string
	host           string
	port           int
	sshClient      *ssh.Client
	*sftp.Client
}

func (client *Client) String() string {
	return fmt.Sprintf("%s@%s:%d", client.username, client.host, client.port)
}

func (client *Client) Close() {
	client.sshClient.Close()
	client.Client.Close()
}

func Connect(host string, port int, certKeyFile string, password string, username string) (*Client, error) {
	certKey, err := isValidFile(certKeyFile)
	defer certKey.Close()
	if err != nil && len(strings.TrimSpace(password)) == 0 {
		return nil, fmt.Errorf("private key[%s] is invalid or password[%s] is empty, inner err[%w]", certKeyFile, password, err)
	}
	config, err := NewClientConfig(5*time.Second, username, certKey, strings.TrimSpace(password))
	if err != nil {
		return nil, err
	}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	addrs := strings.Split(sshClient.RemoteAddr().String(), ":")
	if len(addrs) != 2 {
		return nil, errors.New(fmt.Sprintf("the host[%s] is not standard ", sshClient.RemoteAddr()))
	}
	port, err = strconv.Atoi(addrs[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("the port of host[%s] is not integer", sshClient.RemoteAddr()))
	}
	return &Client{
		username:    username,
		certKeyFile: certKeyFile,
		password:    password,
		host:        addrs[0],
		port:        port,
		sshClient:   sshClient,
		Client:      sftpClient,
	}, nil
}

func ConnectKnownHost(host string, port int, certKeyFile string, knownHostsFile string, username string) (*Client, error) {
	certKey, err := isValidFile(certKeyFile)
	defer certKey.Close()
	if err != nil {
		return nil, fmt.Errorf("private key[%s] is invalid, err[%w]", certKeyFile, err)
	}
	knownHosts, err := isValidFile(knownHostsFile)
	defer knownHosts.Close()
	if err != nil {
		return nil, fmt.Errorf("known_hosts[%s] is invalid, err[%w]", knownHostsFile, err)
	}
	config, err := NewClientConfigWithinKnownHosts(5*time.Second, username, certKey, knownHosts, host)
	if err != nil {
		return nil, err
	}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		log.Fatalf("ssh client dial error : %v", err)
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatalf("new sftp client error : %v", err)
		return nil, err
	}

	addrs := strings.Split(sshClient.RemoteAddr().String(), ":")
	if len(addrs) != 2 {
		return nil, errors.New(fmt.Sprintf("the host[%s] is not standard ", sshClient.RemoteAddr()))
	}
	port, err = strconv.Atoi(addrs[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("the port of host[%s] is not integer", sshClient.RemoteAddr()))
	}
	return &Client{
		username:       username,
		certKeyFile:    certKeyFile,
		knownHostsFile: knownHostsFile,
		host:           addrs[0],
		port:           port,
		sshClient:      sshClient,
		Client:         sftpClient,
	}, nil
}

func isValidFile(targetFile string) (*os.File, error) {
	if certKeyFileInfo, certKeyErr := os.Stat(targetFile); certKeyErr != nil {
		if os.IsNotExist(certKeyErr) {
			certKeyErr = fmt.Errorf("%s is not exist, err[%w]", certKeyFileInfo, certKeyErr)
		}
		return nil, certKeyErr
	} else if certKeyFileInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is directory", targetFile))
	}
	if certFile, err := os.Open(targetFile); err != nil {
		return nil, err
	} else {
		return certFile, nil
	}
}

func NewClientConfig(timeout time.Duration, username string, certKey io.Reader, password string) (*ssh.ClientConfig, error) {
	authMethod := make([]ssh.AuthMethod, 0, 4)
	if certKey != nil {
		key, err := io.ReadAll(certKey)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		authMethod = append(authMethod, ssh.PublicKeys(signer))
	} else if len(strings.TrimSpace(password)) != 0 {
		authMethod = append(authMethod, ssh.Password(strings.TrimSpace(password)))
	} else {
		return nil, errors.New("no key or password is specified")
	}
	return &ssh.ClientConfig{
		Timeout:         timeout,
		User:            username,
		Auth:            authMethod,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// NewClientConfigWithinKnownHosts - 返回的ssh.ClientConfig.HostKeyCallback 使用 ssh.FixedHostKey 来校验公钥是否正确。
// 这里使用ssh-keygen 命令生成的公钥和私钥进行登录ssh远程服务器.
// 注意： 这里需要注意一个ssh的配置。配置在 /etc/ssh/ssh_config中.
// 如果ssh_config中的 HashKnownHosts yes, 那么客户端 known_hosts 中的host key的服务器ip地址是经过哈希加密的。
// 如果ssh_config中的 HashKnownHosts no, 那么客户端 known_hosts 中的host key的服务器ip地址就是明文的。
// 因此，如果想使用该方法，需要将HashKnownHosts 置为no, 使用明文, 这样才能匹配到对应的host key 的ip地址。
func NewClientConfigWithinKnownHosts(timeout time.Duration, username string, certKey io.Reader, knownHosts io.Reader, host string) (*ssh.ClientConfig, error) {
	authMethod := make([]ssh.AuthMethod, 0, 4)
	key, err := io.ReadAll(certKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	authMethod = append(authMethod, ssh.PublicKeys(signer))

	return &ssh.ClientConfig{
		Timeout: timeout,
		User:    username,
		Auth:    authMethod,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			var hostKey ssh.PublicKey
			scanner := bufio.NewScanner(knownHosts)
			for scanner.Scan() {
				fields := strings.Split(scanner.Text(), " ")
				if len(fields) != 3 {
					continue
				}
				if !strings.EqualFold(fields[1], key.Type()) {
					continue
				}
				if strings.Contains(fields[0], host) {
					var err error
					hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
					if err != nil {
						return fmt.Errorf("error parsing %q: %w", fields[2], err)
					}
					break
				}
			}
			if hostKey == nil {
				return errors.New(fmt.Sprintf("no host key for %s", host))
			}
			return nil
		},
	}, nil
}
