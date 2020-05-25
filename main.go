package main // import "github.com/adisbladis/activated-tunnel"

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

func main() {
	serverEndpoint := &Endpoint{
		Host: "159.69.86.193",
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: "localhost",
		Port: 8080,
	}

	hostKey, err := getHostKey(serverEndpoint.String())
	if err != nil {
		panic(err)
	}

	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if bytes.Compare(key.Marshal(), hostKey.Marshal()) == 0 {
			return nil
		}
		return fmt.Errorf("Host key mismatch")
	}

	sshConfig := &ssh.ClientConfig{
		User: "adisbladis",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: []string{hostKey.Type()},
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
