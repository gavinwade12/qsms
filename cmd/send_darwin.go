// +build darwin

package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/nyaruka/phonenumbers"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an sms",
	Long: `Sending an sms requires a carrier, number, and text:
qsms send -n 41955512345 -t "Hello World!"`,
	PreRun: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			exitWithError(err)
		}

		filename := path.Join(home, ".imessage.sh")
		_, err = os.Stat(filename)
		if err != nil && !os.IsNotExist(err) {
			exitWithError(err)
		}
		if err == nil {
			return
		}

		f, err := os.Create(filename)
		if err != nil {
			exitWithError(err)
		}
		defer f.Close()

		f.WriteString(`#!/bin/sh
if [ "$#" -eq 1 ]; then stdinmsg=$(cat); fi
exec <"$0" || exit; read v; read v; read v; exec /usr/bin/osascript - "$@" "$stdinmsg"; exit

on run {phoneNumber, message}
	activate application "Messages"
	tell application "System Events" to tell process "Messages"
		key code 45 using command down -- press Command + N to start a new window
		keystroke phoneNumber -- input the phone number
		key code 36 -- press Enter to focus on the message area 
		keystroke message -- type some message
		key code 36 -- press Enter to send
	end tell
	tell application "Messages" to close window 1 -- Messages was likely not open since qsms was used, so close the window
end run`)
	},
	Run: func(cmd *cobra.Command, args []string) {
		number, err := cmd.Flags().GetString("number")
		if err != nil {
			exitWithError(err)
		}

		text, err := cmd.Flags().GetString("text")
		if err != nil {
			exitWithError(err)
		}

		num, err := phonenumbers.Parse(number, "US")
		if err != nil {
			exitWithError(err)
		}
		number = phonenumbers.Format(num, phonenumbers.E164)[2:]

		home, err := homedir.Dir()
		if err != nil {
			exitWithError(err)
		}

		c := exec.Command("/bin/sh", path.Join(home, ".imessage.sh"), number, text)
		if err = c.Run(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("number", "n", "", "the recipient phone number")
	sendCmd.MarkFlagRequired("number")
	sendCmd.Flags().StringP("text", "t", "", "the sms body")
	sendCmd.MarkFlagRequired("text")
}
