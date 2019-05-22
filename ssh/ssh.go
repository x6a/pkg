// Copyright (C) 2019 <x6a@7n.io>
//
// pkg is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// pkg is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with pkg. If not, see <http://www.gnu.org/licenses/>.

package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type sshCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type sshClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

func (client *sshClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,  // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func (client *sshClient) prepareCommand(session *ssh.Session, cmd *sshCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

func (client *sshClient) runCommand(cmd *sshCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func publicKeyAuthMethod(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// ExecSSH executes a command via SSH on a remote host
func ExecSSH(username, sshPrivateKeyFile, host string, port int, command string) error {

	timeout, _ := time.ParseDuration("20s")

	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{publicKeyAuthMethod(sshPrivateKeyFile)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	client := &sshClient{
		Config: sshConfig,
		Host:   host,
		Port:   port,
	}

	cmd := &sshCommand{
		//Path:   "ls -l $LC_DIR",
		Path:   command,
		Env:    []string{"LC_DIR=/tmp"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := client.runCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}

	return nil
}
