package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
	"os/user"
)

func getHostKey(hostname string) (ssh.PublicKey, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	kh, err := os.Open(usr.HomeDir + "/.ssh/known_hosts")
	if err != nil {
		fmt.Println("unable to read known hosts: %v", err)
		return nil, err
	}

	hx := knownhosts.Normalize(hostname)

	scanner := bufio.NewScanner(kh)
	line := 0
	for scanner.Scan() {
		line += 1
		if len(scanner.Bytes()) == 0 {
			continue
		}
		_, hosts, pubKey, _, _, err := ssh.ParseKnownHosts(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		for _, h := range hosts {
			if h == hx {
				return pubKey, nil
			}
		}
	}

	return nil, fmt.Errorf("No public key found for host")
}
