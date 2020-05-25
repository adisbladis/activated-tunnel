package main

import (
	"fmt"
	"github.com/armon/go-socks5"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/net/context"
	"io"
	"net"
	"os"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Server    *Endpoint
	Remote    *Endpoint
	Config    *ssh.ClientConfig
	Forwarder func(*ssh.Client, net.Conn)
}

func (tunnel *SSHtunnel) Start() error {

	listeners, err := ListenSystemdFds()
	if err != nil {
		panic(err)
	}

	if len(listeners) < 1 {
		panic("Unexpected number of socket activation fds")
	}

	connections := make(chan net.Conn)

	for _, listener := range listeners {
		go func(l net.Listener) {
			for {
				c, err := l.Accept()
				if err != nil {
					panic(err)
					return
				}
				connections <- c
			}
		}(listener)
	}

	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case c := <-connections:
			go func() {
				tunnel.Forwarder(serverConn, c)
			}()
		}
	}

	return nil
}

func (tunnel *SSHtunnel) forwardSocks(serverConn *ssh.Client, localConn net.Conn) {

	dial := func(ctx context.Context, net_, addr string) (net.Conn, error) {
		return serverConn.Dial(net_, addr)
	}

	conf := &socks5.Config{
		Dial: dial,
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	server.ServeConn(localConn)
}

func (tunnel *SSHtunnel) forwardPort(serverConn *ssh.Client, localConn net.Conn) {
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
