package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	InitConfig("config.json")
	to := "ribencong@126.com"
	body := "Hello, this is a test email from Go."

	err := sendEmail([]string{to}, body)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func TestCreateDefaultConfigFile(t *testing.T) {
	cfg := &CliConfig{
		SmtpHost:       "smtp.163.com",
		SmtpPort:       465,
		SenderEmail:    "ribencong@163.com",
		SenderPassword: "[PASSWORD]",
	}

	bts, _ := json.MarshalIndent(cfg, "", "\t")
	_ = os.WriteFile("config.sample.json", bts, 0644)
}
