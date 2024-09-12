package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/realbmail/ex_sdk_example/common"
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"strings"
)

type ClitPara struct {
	to       string
	msg      string
	password string
	config   string
	version  bool
	simple   bool
}

var (
	param = &ClitPara{}
)

var rootCmd = &cobra.Command{
	Use: "bmailCli",

	Short: "bmailCli",

	Long: `usage description::TODO::`,

	Run: mainRun,
}

func init() {
	flags := rootCmd.Flags()
	flags.BoolVarP(&param.simple, "simple",
		"s", false, "bmailCli.lnx -s")
	flags.BoolVarP(&param.version, "version",
		"v", false, "bmailCli.lnx -v")
	flags.StringVarP(&param.to, "to",
		"t", "", "bmailCli.lnx -t ")
	flags.StringVarP(&param.msg, "message",
		"m", "", "bmailCli.lnx -m")
	flags.StringVarP(&param.password, "password",
		"p", "", "bmailCli.lnx -p ")
	flags.StringVarP(&param.config, "config",
		"c", "config.json", "bmailCli.lnx -c ")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

type SdkParam struct {
	Emails []string `json:"emails"`
	Msg    string   `json:"msg"`
}

type SdkResult struct {
	Success      bool              `json:"success"`
	ErrCode      string            `json:"err_code"`
	BMail        map[string]string `json:"b_mail"`
	EncryptedMsg string            `json:"encrypted_msg"`
}

var api = "/encrypt_mail"

func mainRun(_ *cobra.Command, _ []string) {
	if param.version {
		common.PrintVersion()
		return
	}

	InitConfig(param.config)
	receivers := strings.Split(param.to, ",")
	if param.simple {
		fmt.Println("result:=>", sendEmail(receivers, param.msg))
		return
	}
	req := &SdkParam{
		Emails: strings.Split(param.to, ","),
		Msg:    param.msg,
	}
	reqData, _ := json.Marshal(req)

	respData, err := doHttp(_cliCfg.Server+api, "application/json", reqData)
	if err != nil {
		fmt.Println("failed to encrypt mail data:", err)
		return
	}
	var result = &SdkResult{}
	err = json.Unmarshal(respData, result)
	if err != nil {
		fmt.Println("failed to encrypt mail data:", err)
		return
	}

	if !result.Success {
		fmt.Println("failed to encrypt mail data:=>", result.EncryptedMsg)
		return
	}

	fmt.Println("according bmail address=>", result.BMail)
	_ = sendEmail(receivers, result.EncryptedMsg)
}

func sendEmail(tos []string, body string) error {
	smtpHost := _cliCfg.SmtpHost
	smtpPort := _cliCfg.SmtpPort
	senderEmail := _cliCfg.SenderEmail

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	for _, to := range tos {
		m.SetHeader("To", to)
	}

	m.SetHeader("Subject", "Test Email from Go")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, _cliCfg.SenderPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("Send Mail Success")
	return nil
}

func doHttp(url, cTyp string, data []byte) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", cTyp)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
