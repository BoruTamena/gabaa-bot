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
	"go.uber.org/zap"
)

type Client struct {
	apiKey     string
	baseURL    string
	publicKey  string
	httpClient *http.Client
}

func NewClient() *Client {
	cfg, err := loadClientConfig()
	if err != nil {
		logger.Error("lakipay client not configured", zap.Error(err))
		return &Client{
			httpClient: &http.Client{Timeout: 30 * time.Second},
		}
	}

	return &Client{
		apiKey:     cfg.apiKey,
		baseURL:    cfg.baseURL,
		publicKey:  cfg.webhookPublicKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type clientConfig struct {
	apiKey           string
	webhookPublicKey string
	baseURL          string
}

func loadClientConfig() (*clientConfig, error) {
	baseURL := configString("lakipay.base.url", "lakipay.baseurl")
	if baseURL == "" {
		baseURL = "https://api.lakipay.co"
	}

	secretKey := configString("lakipay.secret.key", "lakipay.secretkey")
	pubKey := configString("lakipay.pub.key", "lakipay.api.secret")

	webhookPublicKey := configString("lakipay.public.key", "lakipay.webhook_public_key")
	if webhookPublicKey != "" && !strings.Contains(webhookPublicKey, "BEGIN PUBLIC KEY") {
		// Support older setups that stored the API secret in LAKIPAY_PUBLIC_KEY.
		// if apiSecret == "" {
		// 	apiSecret = webhookPublicKey
		// }
		webhookPublicKey = ""
	}

	// if secretKey == "" {
	// 	return nil, fmt.Errorf("set LAKIPAY_SECRET_KEY")
	// }
	// if apiSecret == "" {
	// 	return nil, fmt.Errorf("set LAKIPAY_PUB_KEY (LAKISEC_...) for API authentication")
	// }

	return &clientConfig{
		apiKey:           pubKey + ":" + secretKey,
		webhookPublicKey: webhookPublicKey,
		baseURL:          strings.TrimRight(baseURL, "/"),
	}, nil
}

func configString(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(viper.GetString(key)); value != "" {
			return value
		}
	}
	return ""
}

func (c *Client) InitiateDirectPayment(ctx context.Context, req DirectPaymentRequest) (*DirectPaymentResponse, error) {
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("lakipay client not configured")
	}
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
		logger.Error("lakipay request failed", zap.Error(err), zap.String("api_key", c.apiKey))
		return nil, fmt.Errorf("lakipay request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("lakipay request failed", zap.Error(err), zap.String("api_key", c.apiKey))
		return nil, fmt.Errorf("read response: %w", err)
	}

	logger.Info("lakipay response", zap.String("response", string(respBody)), zap.String("api_key", c.apiKey))
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
	if c == nil || c.apiKey == "" {
		return nil, fmt.Errorf("lakipay client not configured")
	}
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

// ConfigurationError reports whether required LakiPay credentials are missing.
func (c *Client) ConfigurationError() error {
	if c == nil || c.apiKey == "" {
		_, err := loadClientConfig()
		return err
	}
	return nil
}

func (c *Client) VerifyWebhookSignature(payload map[string]string, signature string) (bool, error) {
	if c == nil || c.publicKey == "" {
		return false, fmt.Errorf("lakipay webhook public key not configured")
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
