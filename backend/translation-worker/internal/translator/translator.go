package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/woragis/backend/translation-worker/internal/config"
)

// Translator provides translation services via external APIs.
type Translator interface {
	Translate(ctx context.Context, text string, targetLanguage string) (string, error)
}

// NewTranslator creates a translator based on the configured provider.
func NewTranslator(cfg config.TranslationConfig, logger *slog.Logger) (Translator, error) {
	switch cfg.Provider {
	case config.ProviderGoogle:
		if cfg.GoogleAPIKey == "" {
			return nil, fmt.Errorf("GOOGLE_TRANSLATE_API_KEY is required for Google Translate")
		}
		return NewGoogleTranslator(cfg.GoogleAPIKey, cfg.GoogleProjectID, cfg.Timeout, cfg.MaxRetries, cfg.RetryDelay, logger), nil
	case config.ProviderDeepL:
		if cfg.DeepLAPIKey == "" {
			return nil, fmt.Errorf("DEEPL_API_KEY is required for DeepL")
		}
		return NewDeepLTranslator(cfg.DeepLAPIKey, cfg.Timeout, cfg.MaxRetries, cfg.RetryDelay, logger), nil
	case config.ProviderLibre:
		apiURL := cfg.LibreAPIURL
		if apiURL == "" {
			apiURL = "https://libretranslate.com/translate" // Default public instance
		}
		return NewLibreTranslator(apiURL, cfg.LibreAPIKey, cfg.Timeout, cfg.MaxRetries, cfg.RetryDelay, logger), nil
	default:
		return nil, fmt.Errorf("unsupported translation provider: %s", cfg.Provider)
	}
}

// GoogleTranslator implements translation using Google Cloud Translation API.
type GoogleTranslator struct {
	apiKey     string
	projectID  string
	timeout    int
	maxRetries int
	retryDelay int
	logger     *slog.Logger
	client     *http.Client
}

// NewGoogleTranslator creates a new Google Translate client.
func NewGoogleTranslator(apiKey, projectID string, timeout, maxRetries, retryDelay int, logger *slog.Logger) *GoogleTranslator {
	return &GoogleTranslator{
		apiKey:     apiKey,
		projectID:  projectID,
		timeout:    timeout,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		logger:     logger,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Translate translates text using Google Cloud Translation API.
func (t *GoogleTranslator) Translate(ctx context.Context, text string, targetLanguage string) (string, error) {
	// Map language codes to Google's format
	targetLang := mapLanguageToGoogle(targetLanguage)

	// Google Cloud Translation API v3 endpoint
	apiURL := "https://translation.googleapis.com/v3/projects"
	if t.projectID != "" {
		apiURL = fmt.Sprintf("%s/%s", apiURL, t.projectID)
	}
	apiURL = fmt.Sprintf("%s:translateText", apiURL)

	requestBody := map[string]interface{}{
		"contents":           []string{text},
		"targetLanguageCode": targetLang,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < t.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(t.retryDelay) * time.Millisecond)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))

		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			t.logger.Warn("Translation request failed, retrying",
				slog.Int("attempt", attempt+1),
				slog.Any("error", err),
			)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				// Don't retry client errors
				return "", lastErr
			}
			t.logger.Warn("Translation API error, retrying",
				slog.Int("attempt", attempt+1),
				slog.Int("status", resp.StatusCode),
			)
			continue
		}

		var result struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("failed to decode response: %w", err)
			continue
		}

		if len(result.Translations) == 0 {
			return "", fmt.Errorf("no translation in response")
		}

		return strings.TrimSpace(result.Translations[0].TranslatedText), nil
	}

	return "", fmt.Errorf("translation failed after %d attempts: %w", t.maxRetries, lastErr)
}

// DeepLTranslator implements translation using DeepL API.
type DeepLTranslator struct {
	apiKey     string
	timeout    int
	maxRetries int
	retryDelay int
	logger     *slog.Logger
	client     *http.Client
}

// NewDeepLTranslator creates a new DeepL client.
func NewDeepLTranslator(apiKey string, timeout, maxRetries, retryDelay int, logger *slog.Logger) *DeepLTranslator {
	return &DeepLTranslator{
		apiKey:     apiKey,
		timeout:    timeout,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		logger:     logger,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Translate translates text using DeepL API.
func (t *DeepLTranslator) Translate(ctx context.Context, text string, targetLanguage string) (string, error) {
	// Map language codes to DeepL's format
	targetLang := mapLanguageToDeepL(targetLanguage)

	// Determine API endpoint (free or pro)
	apiURL := "https://api-free.deepl.com/v2/translate"
	if strings.HasPrefix(t.apiKey, "DeepL-Auth-Key") {
		apiURL = "https://api.deepl.com/v2/translate"
	}

	formData := url.Values{}
	formData.Set("auth_key", t.apiKey)
	formData.Set("text", text)
	formData.Set("target_lang", targetLang)

	var lastErr error
	for attempt := 0; attempt < t.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(t.retryDelay) * time.Millisecond)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			t.logger.Warn("Translation request failed, retrying",
				slog.Int("attempt", attempt+1),
				slog.Any("error", err),
			)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", lastErr
			}
			t.logger.Warn("Translation API error, retrying",
				slog.Int("attempt", attempt+1),
				slog.Int("status", resp.StatusCode),
			)
			continue
		}

		var result struct {
			Translations []struct {
				Text string `json:"text"`
			} `json:"translations"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("failed to decode response: %w", err)
			continue
		}

		if len(result.Translations) == 0 {
			return "", fmt.Errorf("no translation in response")
		}

		return strings.TrimSpace(result.Translations[0].Text), nil
	}

	return "", fmt.Errorf("translation failed after %d attempts: %w", t.maxRetries, lastErr)
}

// LibreTranslator implements translation using LibreTranslate API.
type LibreTranslator struct {
	apiURL     string
	apiKey     string
	timeout    int
	maxRetries int
	retryDelay int
	logger     *slog.Logger
	client     *http.Client
}

// NewLibreTranslator creates a new LibreTranslate client.
func NewLibreTranslator(apiURL, apiKey string, timeout, maxRetries, retryDelay int, logger *slog.Logger) *LibreTranslator {
	return &LibreTranslator{
		apiURL:     apiURL,
		apiKey:     apiKey,
		timeout:    timeout,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		logger:     logger,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Translate translates text using LibreTranslate API.
func (t *LibreTranslator) Translate(ctx context.Context, text string, targetLanguage string) (string, error) {
	// Map language codes to LibreTranslate's format
	targetLang := mapLanguageToLibre(targetLanguage)

	requestBody := map[string]interface{}{
		"q":      text,
		"source": "auto", // Auto-detect source language
		"target": targetLang,
		"format": "text",
	}

	if t.apiKey != "" {
		requestBody["api_key"] = t.apiKey
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < t.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(t.retryDelay) * time.Millisecond)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", t.apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			t.logger.Warn("Translation request failed, retrying",
				slog.Int("attempt", attempt+1),
				slog.Any("error", err),
			)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", lastErr
			}
			t.logger.Warn("Translation API error, retrying",
				slog.Int("attempt", attempt+1),
				slog.Int("status", resp.StatusCode),
			)
			continue
		}

		var result struct {
			TranslatedText string `json:"translatedText"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = fmt.Errorf("failed to decode response: %w", err)
			continue
		}

		return strings.TrimSpace(result.TranslatedText), nil
	}

	return "", fmt.Errorf("translation failed after %d attempts: %w", t.maxRetries, lastErr)
}

// Language mapping functions
func mapLanguageToGoogle(lang string) string {
	langMap := map[string]string{
		"en":    "en",
		"pt-BR": "pt",
		"pt":    "pt",
		"es":    "es",
		"fr":    "fr",
		"de":    "de",
		"ru":    "ru",
		"ja":    "ja",
		"ko":    "ko",
		"zh-CN": "zh",
		"zh":    "zh",
		"el":    "el",
		"la":    "la",
		"sv":    "sv",
	}
	if mapped, ok := langMap[lang]; ok {
		return mapped
	}
	return lang
}

func mapLanguageToDeepL(lang string) string {
	langMap := map[string]string{
		"en":    "EN",
		"pt-BR": "PT-BR",
		"pt":    "PT",
		"es":    "ES",
		"fr":    "FR",
		"de":    "DE",
		"ru":    "RU",
		"ja":    "JA",
		"ko":    "KO",
		"zh-CN": "ZH",
		"zh":    "ZH",
		"el":    "EL",
		"la":    "LA",
		"sv":    "SV",
	}
	if mapped, ok := langMap[lang]; ok {
		return mapped
	}
	return strings.ToUpper(lang)
}

func mapLanguageToLibre(lang string) string {
	langMap := map[string]string{
		"en":    "en",
		"pt-BR": "pt",
		"pt":    "pt",
		"es":    "es",
		"fr":    "fr",
		"de":    "de",
		"ru":    "ru",
		"ja":    "ja",
		"ko":    "ko",
		"zh-CN": "zh",
		"zh":    "zh",
		"el":    "el",
		"la":    "la",
		"sv":    "sv",
	}
	if mapped, ok := langMap[lang]; ok {
		return mapped
	}
	return lang
}
