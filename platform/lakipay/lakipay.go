package lakipay

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/spf13/viper"
)

type Client struct {
	apiKey     string
	baseURL    string
	publicKey  string
	httpClient *http.Client
}

func NewClient() *Client {
	baseURL := viper.GetString("lakipay.base_url")
	if baseURL == "" {
		baseURL = "https://api.lakipay.co"
	}

	secretKey := viper.GetString("lakipay.secret_key")
	if secretKey == "" {
		logger.Error("lakipay secret key not configured")
		return nil
	}
	pubKey := viper.GetString("lakipay.public_key")
	if pubKey == "" {
		logger.Error("lakipay public key not configured")
		return nil
	}

	apikey := secretKey + ":" + pubKey
	return &Client{
		apiKey:    apikey,
		baseURL:   strings.TrimRight(baseURL, "/"),
		publicKey: viper.GetString("lakipay.public_key"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) InitiateDirectPayment(ctx context.Context, req DirectPaymentRequest) (*DirectPaymentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v2/payment/direct", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("lakipay request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lakipay error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result DirectPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if strings.ToUpper(result.Status) != "SUCCESS" {
		return nil, fmt.Errorf("lakipay rejected payment: %s", result.Message)
	}

	return &result, nil
}

func (c *Client) InitiateWithdrawal(ctx context.Context, req WithdrawalRequest) (*WithdrawalResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v2/payment/withdrawal", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("lakipay withdrawal request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lakipay error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result WithdrawalResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if strings.ToUpper(result.Status) != "SUCCESS" {
		return nil, fmt.Errorf("lakipay rejected withdrawal: %s", result.Message)
	}

	return &result, nil
}

func (c *Client) VerifyWebhookSignature(payload map[string]string, signature string) (bool, error) {
	if c.publicKey == "" {
		return false, fmt.Errorf("lakipay public key not configured")
	}
	if signature == "" {
		return false, fmt.Errorf("missing signature")
	}

	canonical := buildCanonicalString(payload)
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("decode signature: %w", err)
	}

	pubKey, err := parsePublicKey(c.publicKey)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256([]byte(canonical))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func buildCanonicalString(payload map[string]string) string {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, payload[k]))
	}
	return strings.Join(pairs, "&")
}

func parsePublicKey(pemKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format
		rsaPub, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse public key: %w", err)
		}
		return rsaPub, nil
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPub, nil
}

// PayloadToStringMap converts a webhook JSON body to a flat string map for signature verification.
func PayloadToStringMap(raw json.RawMessage) (map[string]string, error) {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		return nil, err
	}

	result := make(map[string]string, len(rawMap))
	for k, v := range rawMap {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			// Use raw JSON value without quotes for numbers/bools
			result[k] = strings.Trim(string(v), `"`)
		} else {
			result[k] = s
		}
	}
	return result, nil
}
