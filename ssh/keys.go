package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func makeKeyring() []ssh.Signer {
	signers := []ssh.Signer{}
	// TODO(ChrisMcKenzie): makekeyring from paths in env vars
	keys := []string{
		os.Getenv("KEY_PATH") + "/.ssh/id_rsa",
		os.Getenv("KEY_PATH") + "/.ssh/id_dsa",
		os.Getenv("PUB_KEY_RSA"),
	}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return signers
}

func NewClientConfig(u string, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: u,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(makeKeyring()...),
			ssh.Password(password),
		},
	}
}

func Execute(commands []string, server string, config *ssh.ClientConfig) string {
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return "Failed to dial: " + err.Error()
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "Failed to create session: " + err.Error()
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b

	script := "set -x\n" + strings.Join(commands, ";set +x;echo \"--------------\"; set -x;")

	session.Run(script)

	return fmt.Sprintf("[%s]: %s", server, b.String())
}
