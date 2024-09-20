package models

type ShortenDto struct {
	URL string `json:"url"`
}

type ShortDto struct {
	Result string `json:"result"`
}

type ShortenBatchItemDto struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortBatchItemDto struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserUrlDto struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
