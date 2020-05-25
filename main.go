package main // import "github.com/adisbladis/activated-tunnel"

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"os/user"
)

func main() {

	var username string

	tunnel := &SSHtunnel{
		Server: &Endpoint{},
		Remote: &Endpoint{},
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "Server address",
				Required:    true,
				Destination: &tunnel.Server.Host,
			},
			&cli.StringFlag{
				Name:        "user",
				Usage:       "Server username",
				Value:       usr.Username,
				Destination: &username,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       22,
				Usage:       "Server port",
				Destination: &tunnel.Server.Port,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "port",
				Usage: "Forward a single port on the remote (ssh -L)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "host",
						Usage:       "Remote host",
						Value:       "localhost",
						Destination: &tunnel.Remote.Host,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "Remote port",
						Required:    true,
						Destination: &tunnel.Remote.Port,
					},
				},
				Action: func(c *cli.Context) error {
					tunnel.Forwarder = tunnel.forwardPort
					return nil
				},
			},
			{
				Name:  "socks",
				Usage: "Run in SOCKS mode (ssh -D)",
				Action: func(c *cli.Context) error {
					tunnel.Forwarder = tunnel.forwardSocks
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	hostKey, err := getHostKey(tunnel.Server.String())
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
		User: username,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: []string{hostKey.Type()},
	}

	tunnel.Config = sshConfig

	tunnel.Start()
}
