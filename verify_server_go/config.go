package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CliConfig struct {
	SmtpHost       string `json:"smtp_host"`
	SmtpPort       int    `json:"smtp_port"`
	SenderEmail    string `json:"sender_email"`
	SenderPassword string `json:"sender_password"`
	Server         string `json:"server"`
}

var _cliCfg *CliConfig

func (c *CliConfig) String() string {
	s := "\n------system config------"
	s += "\nsmtp server:\t" + c.SmtpHost
	s += "\nsmtp host:\t" + fmt.Sprintf("%d", c.SmtpPort)
	s += "\nsmtp sender:\t" + c.SenderEmail
	s += "\nsdk server:\t" + c.Server
	s += "\n-------------------------"
	return s
}

func InitConfig(filName string) *CliConfig {
	cf := new(CliConfig)

	bts, err := os.ReadFile(filName)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(bts, &cf); err != nil {
		panic(err)
	}

	_cliCfg = cf
	fmt.Println(cf.String())
	return cf
}
