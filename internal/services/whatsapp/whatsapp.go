package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/juniorAkp/easyPay/pkg/types"
)

type TemplateName string

const (
	HelloWorld   TemplateName = "hello_world"
	EasyPayIntro TemplateName = "easy_pay_intro"
)

const whatsappAPI = "https://graph.facebook.com/v21.0"

func SendMessage(to, message string) error {
	data := types.MessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
	}
	data.Text.Body = message

	jsonData, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/%s/messages", whatsappAPI, os.Getenv("PHONE_ID"))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("WHATSAPP_ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

func SendTemplateMessage(to, templateName string) error {
	data := types.TemplateMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
	}
	data.Template.Name = templateName
	data.Template.Language.Code = "en_US"

	jsonData, _ := json.Marshal(data)

	url := fmt.Sprintf("%s/%s/messages", whatsappAPI, os.Getenv("PHONE_ID"))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("WHATSAPP_ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response", string(body))
	return nil
}
