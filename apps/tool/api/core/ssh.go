package core

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type SSHConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func NewSshClient(scfg SSHConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            scfg.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		Auth:            []ssh.AuthMethod{ssh.Password(scfg.Password)},
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", scfg.Host, scfg.Port), config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func hostKeyCallBackFunc(host string) ssh.HostKeyCallback {
	hostPath, err := homedir.Expand("~/.ssh/known_hosts")
	if err != nil {
		log.Fatal("find known_hosts's home dir failed", err)
	}
	file, err := os.Open(hostPath)
	if err != nil {
		log.Fatal("can't find known_host file:", err)
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
		log.Fatalf("no hostkey for %s,%v", host, err)
	}
	return ssh.FixedHostKey(hostKey)
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}
