package types

type MessageRequest struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Field string `json:"field"`
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

type MessageDetails struct {
	Account string  `json:"account"`
	Amount  float64 `json:"amount"`
	Message string  `json:"message"`
}

type TemplateMessage struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Template         struct {
		Name     string `json:"name"`
		Language struct {
			Code string `json:"code"`
		} `json:"language"`
	} `json:"template"`
}

type RequestToPayRequest struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalId   string `json:"externalId"`
	Payer        Payer  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

type Payer struct {
	PartyIdType string `json:"partyIdType"`
	PartyId     string `json:"partyId"`
}

type IncomingMessage struct {
	Phone     string  `json:"phone"`
	Name      string  `json:"name"`
	Message   string  `json:"message"`
	Timestamp string  `json:"timestamp"`
	Account   string  `json:"account"`
	Amount    float64 `json:"amount"`
}
