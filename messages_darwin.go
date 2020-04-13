// +build darwin

package main

import (
	"os"
	"os/exec"
	"path"

	"github.com/mitchellh/go-homedir"
)

// MessagesGateway is a Gateway for sending SMS via the OSX Messages app
type MessagesGateway struct{}

// Send implements Gateway.Send
func (g MessagesGateway) Send(recipient, text string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	filename := path.Join(home, ".qsms.messages.sh")

	if err := g.ensureScriptFileExists(filename); err != nil {
		return err
	}

	c := exec.Command("/bin/sh", filename, recipient, text)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func (g MessagesGateway) ensureScriptFileExists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(`#!/bin/sh
if [ "$#" -eq 1 ]; then stdinmsg=$(cat); fi
exec <"$0" || exit; read v; read v; read v; exec /usr/bin/osascript - "$@" "$stdinmsg"; exit

on run {phoneNumber, message}
activate application "Messages"
delay 1
tell application "System Events" to tell process "Messages"
	key code 45 using command down -- press Command + N to start a new window
	keystroke phoneNumber -- input the phone number
	key code 36 -- press Enter to focus on the message area 
	keystroke message -- type some message
	key code 36 -- press Enter to send
end tell
delay 1
tell application "Messages" to close window 1 -- Messages was likely not open since qsms was used, so close the window
end run
`)

	return nil
}
