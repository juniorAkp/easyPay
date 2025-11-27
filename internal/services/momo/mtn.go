package momo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/juniorAkp/easyPay/pkg/types"
)

type MoMoService struct {
	subscriptionKey string
	apiUser         string
	apiKey          string
	environment     string
	baseURL         string
	accessToken     string
	tokenExpiry     time.Time
}

func NewMoMoService() *MoMoService {
	return &MoMoService{
		subscriptionKey: os.Getenv("MTN_SUBSCRIPTION_KEY"),
		apiUser:         os.Getenv("MTN_API_USER"),
		apiKey:          os.Getenv("MTN_API_KEY"),
		environment:     os.Getenv("MTN_ENVIRONMENT"),
		baseURL:         os.Getenv("MTN_BASE_URL"),
	}
}

func (m *MoMoService) GetAccessToken() error {
	if time.Now().Before(m.tokenExpiry) {
		return nil //valid token
	}

	url := "https://sandbox.momodeveloper.mtn.com/collection/token/"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(m.apiUser + ":" + m.apiKey))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.subscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get token: %s", string(body))
	}
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	m.accessToken = result.AccessToken
	m.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return nil

}

func (m *MoMoService) RequestPay(amount, phone, message string) (string, error) {

	if err := m.GetAccessToken(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", m.baseURL)

	referenceId := uuid.New().String()

	payload := types.RequestToPayRequest{
		Amount:     amount,
		Currency:   "EUR", //change in prod
		ExternalId: referenceId,
		Payer: types.Payer{
			PartyIdType: "MSISDN",
			PartyId:     phone,
		},
		PayerMessage: message,
		PayeeNote:    "Payment via EasyPay",
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	//headers for mtn
	req.Header.Set("Authorization", "Bearer "+m.accessToken)
	req.Header.Set("X-Reference-Id", referenceId)
	req.Header.Set("X-Target-Environment", m.environment)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Ocp-Apim-Subscription-Key", m.subscriptionKey)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("request failed: %s", string(body))
	}

	return referenceId, nil
}

func (m *MoMoService) CheckTransactionStatus(referenceId string) (string, error) {
	if err := m.GetAccessToken(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", m.baseURL, referenceId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+m.accessToken)
	req.Header.Set("X-Target-Environment", m.environment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.subscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
