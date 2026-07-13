package lakipay

import "strings"

func firstMapValue(m map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(m[key]); value != "" {
			return value
		}
	}
	return ""
}

// WebhookTransactionID extracts the LakiPay transaction id from a webhook payload.
func WebhookTransactionID(m map[string]string) string {
	return firstMapValue(m, "lakipayTxnId", "lakipay_transaction_id", "transaction_id")
}

// WebhookReference extracts the merchant reference from a webhook payload.
func WebhookReference(m map[string]string) string {
	return firstMapValue(m, "referenceId", "reference_id", "reference")
}
