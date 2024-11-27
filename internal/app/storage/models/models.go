package models

type BatchLinkRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchLinkResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserBatchLink struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
