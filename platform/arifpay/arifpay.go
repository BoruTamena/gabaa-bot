package arifpay

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BoruTamena/gabaa-bot/platform"
)

type arifPayment struct {
	apiKey     string
	url        string
	devMode    bool
	sandBoxUrl string
	expireDate time.Time
}

func NewPayment() platform.Payment {

	return &arifPayment{
		apiKey:     "",
		url:        "",
		devMode:    false,
		sandBoxUrl: "",
	}

}

func (ap *arifPayment) MakePayment(payload platform.PaymentRequestPayload) (error, platform.PaymentResponse) {

	var response platform.PaymentResponse
	client := http.Client{} // creating request clients

	body, err := json.Marshal(payload)

	if err != nil {
		log.Println("marshalling request payload err::", err)
		return err, response
	}

	// creating request
	req, err := http.NewRequest(http.MethodPost, ap.url, bytes.NewBuffer(body))

	if err != nil {
		log.Println("creating request err::", err)
		return err, response
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-arifpay-Key", ap.apiKey)
	// making request
	res, err := client.Do(req)

	if err != nil {

		log.Println("making request err::", err)
		return err, response
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {

		log.Println("response body decode err::", err)

		return err, response

	}

	return nil, response

}

func (ap *arifPayment) Transaction() error {

	return nil
}
