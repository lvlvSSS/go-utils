package sftp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
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

func (client *Client) Host() string {
	return client.host
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

	host, p, err := net.SplitHostPort(sshClient.RemoteAddr().String())
	if err != nil {
		return nil, fmt.Errorf("try to analyze the host[%s] error[%w]", sshClient.RemoteAddr().String(), err)
	}
	port, err = strconv.Atoi(p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("the port of host[%s] is not integer", sshClient.RemoteAddr()))
	}
	return &Client{
		username:    username,
		certKeyFile: certKeyFile,
		password:    password,
		host:        host,
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

	host, p, err := net.SplitHostPort(sshClient.RemoteAddr().String())
	if err != nil {
		return nil, fmt.Errorf("try to analyze the host[%s] error[%w]", sshClient.RemoteAddr().String(), err)
	}
	port, err = strconv.Atoi(p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("the port of host[%s] is not integer", sshClient.RemoteAddr()))
	}
	return &Client{
		username:       username,
		certKeyFile:    certKeyFile,
		knownHostsFile: knownHostsFile,
		host:           host,
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

// NewClientConfigWithinKnownHosts - 返回的ssh.ClientConfig.HostKeyCallback 使用 自定义的函数 来校验公钥是否正确。
// 这里使用ssh-keygen 命令生成的公钥和私钥进行登录ssh远程服务器.
// 注意： 这里需要注意一个ssh的配置。配置在 /etc/ssh/ssh_config中.
// 如果ssh_config中的 HashKnownHosts yes, 那么客户端 known_hosts 中的host key的服务器 ip地址/主机名 是经过哈希加密的。
// 如果ssh_config中的 HashKnownHosts no, 那么客户端 known_hosts 中的host key的服务器 ip地址/主机名 就是明文的。
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
		// 校验服务端返回的公钥信息是否和known_hosts文件中的公钥信息匹配.
		// known_hosts文件中保存这host key, 每个host key 由3个部分组成.
		// 1. 主机的ip地址/主机名.(ip地址有可能是经过哈希加密的. 是否加密, 由配置文件 /etc/ssh/ssh_config 的 HashKnownHosts 来决定. 如果为yes, 则进行加密. 否则不用加密.)
		// 2. 加密方式
		// 3. base64编码的公钥信息.
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
				if isHostValid(fields[0], host) {
					var err error
					hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
					if err != nil {
						return fmt.Errorf("error parsing %q: %w", fields[2], err)
					}
					break
				}
			}
			if hostKey == nil {
				return errors.New(fmt.Sprintf("no host key in known_hosts for %s", host))
			}
			return nil
		},
	}, nil
}

// isHostValid - 判断要需要远程到的host服务器主机名是否和known_hosts中的服务器主机名匹配
// host - 用户输入要远程的服务器主机名
// hostInKnownHosts - known_hosts中的服务器主机名(有可能经过哈希加密,是否加密取决于HashKnownHosts配置)
func isHostValid(hostInKnownHosts string, host string) bool {
	// 经过哈希加密的主机名, 由3部分组成：
	// 1. sha的协议版本, 一般都是sha1,因此这里一般都是为 ’1‘
	// 2. 给host进行加密的salt的base64编码
	// 3. 加盐(salt)之后, host哈希加密后的base64编码
	// 最后这3部分用字符 '|' 拼接起来. 详见 encodeHash
	ver, salt, hash, err := decodeHash(hostInKnownHosts)
	if err != nil {
		// 解码失败，此时表示 hostInKnownHosts 不是哈希加密后的主机名
		hostInKnownHosts = knownhosts.Normalize(hostInKnownHosts)
		if len(hostInKnownHosts) == 0 {
			return false
		}
		return strings.Contains(hostInKnownHosts, host)
	}
	// 解码成功，表示 hostInKnownHosts 是哈希加密后的主机名
	encoded := hostInKnownHosts
	// 重新拼接，做一次校验
	if got := encodeHash(ver, salt, hash); got != encoded {
		return false
	}
	// 判断sha的加密版本是否为sha1
	if ver != sha1HashType {
		return false
	}
	// 用相同的salt对host进行加密，来进行判断
	if got := hashHost(host, salt); !bytes.Equal(got, hash) {
		return false
	}
	return true
}
