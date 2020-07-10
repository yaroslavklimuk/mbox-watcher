package main

import (
	"os/exec"
	"strings"
)

type MailSender func(msgContent string, recipientEmail string) error

func sendmailSender(msgContent string, recipientEmail string) (err error) {
	cmd := exec.Command("sendmail", recipientEmail)
	cmd.Stdin = strings.NewReader(msgContent)
	err := cmd.Run()
}

func testSender(msgContent string, recipientEmail string) (err error) {

}
