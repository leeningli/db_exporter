package utils

import(
	"fmt"
	"time"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"os"
	"log"
	"net"
)

func SftpConnect(user, pwd, host string, port int) (*sftp.Client, error) {
	var(
		auth []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		sshClient *ssh.Client
		sftpClient *sftp.Client
		err error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(pwd))
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
        	return nil
    	},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func Connect(user, password, host string, port int) (*ssh.Session, error){
	var(
		auth []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		client *ssh.Client
		session *ssh.Session
		err error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
        		return nil
    		},
	}
	//clientConfig.Config.Ciphers = append(clientConfig.Config.Ciphers, "aes128-cbc")
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil{
		return nil, err
	}
	if session, err = client.NewSession(); err != nil{
		return nil, err
	}
	return session, nil
}

func RemoteExec(user, password, host, command string, port int) (*ssh.Session, error){
	session, err := Connect(user, password, host, port)
	if err != nil{
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run(command)
	return session, err
}
