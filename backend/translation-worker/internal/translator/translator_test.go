package translator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"
	"github.com/woragis/backend/translation-worker/internal/config"
)

func TestNewTranslator(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.TranslationConfig
		wantErr bool
	}{
		{
			name: "google translator",
			cfg: config.TranslationConfig{
				Provider:     config.ProviderGoogle,
				GoogleAPIKey: "test-key",
				Timeout:      30,
				MaxRetries:   3,
				RetryDelay:   1000,
			},
			wantErr: false,
		},
		{
			name: "deepl translator",
			cfg: config.TranslationConfig{
				Provider:    config.ProviderDeepL,
				DeepLAPIKey: "test-key",
				Timeout:     30,
				MaxRetries:  3,
				RetryDelay:  1000,
			},
			wantErr: false,
		},
		{
			name: "libre translator",
			cfg: config.TranslationConfig{
				Provider:      config.ProviderLibre,
				LibreAPIURL:   "https://libretranslate.com/translate",
				Timeout:       30,
				MaxRetries:    3,
				RetryDelay:    1000,
			},
			wantErr: false,
		},
		{
			name: "google translator without key",
			cfg: config.TranslationConfig{
				Provider:     config.ProviderGoogle,
				GoogleAPIKey: "",
				Timeout:      30,
				MaxRetries:   3,
				RetryDelay:   1000,
			},
			wantErr: true,
		},
		{
			name: "deepl translator without key",
			cfg: config.TranslationConfig{
				Provider:    config.ProviderDeepL,
				DeepLAPIKey: "",
				Timeout:     30,
				MaxRetries:  3,
				RetryDelay:  1000,
			},
			wantErr: true,
		},
		{
			name: "unsupported provider",
			cfg: config.TranslationConfig{
				Provider: config.TranslationProvider("invalid"),
				Timeout:  30,
				MaxRetries: 3,
				RetryDelay: 1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.Default()
			translator, err := NewTranslator(tt.cfg, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTranslator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && translator == nil {
				t.Error("NewTranslator() returned nil translator")
			}
		})
	}
}

func TestGoogleTranslator_Translate(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %v", r.Method)
		}

		response := map[string]interface{}{
			"translations": []map[string]string{
				{"translatedText": "Hola, este es un mensaje de prueba para traducción."},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create translator with mock server URL
	translator := NewGoogleTranslator("test-key", "test-project", 30, 3, 1000, slog.Default())
	// Override the API URL for testing
	translator.client = &http.Client{}

	// We can't easily test the full Google API without mocking the URL construction
	// This test verifies the translator can be created
	if translator == nil {
		t.Error("NewGoogleTranslator() returned nil")
	}
}

func TestDeepLTranslator_Translate(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %v", r.Method)
		}

		response := map[string]interface{}{
			"translations": []map[string]string{
				{"text": "Bonjour, ceci est un message de test pour la traduction."},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	translator := NewDeepLTranslator("test-key", 30, 3, 1000, slog.Default())
	if translator == nil {
		t.Error("NewDeepLTranslator() returned nil")
	}
}

func TestLibreTranslator_Translate(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %v", r.Method)
		}

		response := map[string]string{
			"translatedText": "Hallo, dies ist eine Testnachricht zur Übersetzung.",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	translator := NewLibreTranslator(server.URL, "", 30, 3, 1000, slog.Default())
	if translator == nil {
		t.Error("NewLibreTranslator() returned nil")
	}

	// Test actual translation with mock server
	ctx := context.Background()
	result, err := translator.Translate(ctx, "Hello, this is a test message for translation.", "de")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if result != "Hallo, dies ist eine Testnachricht zur Übersetzung." {
		t.Errorf("Translate() = %v, want Hallo, dies ist eine Testnachricht zur Übersetzung.", result)
	}
}

func TestLibreTranslator_Translate_Error(t *testing.T) {
	// Create a mock HTTP server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API key required",
		})
	}))
	defer server.Close()

	translator := NewLibreTranslator(server.URL, "", 30, 3, 1000, slog.Default())
	ctx := context.Background()
	_, err := translator.Translate(ctx, "test", "de")
	if err == nil {
		t.Error("Translate() should return error for 400 status")
	}
}

func TestLanguageMapping(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		google   string
		deepl    string
		libre    string
	}{
		{
			name:   "english",
			lang:   "en",
			google: "en",
			deepl:  "EN",
			libre:  "en",
		},
		{
			name:   "portuguese brazil",
			lang:   "pt-BR",
			google: "pt",
			deepl:  "PT-BR",
			libre:  "pt",
		},
		{
			name:   "chinese",
			lang:   "zh-CN",
			google: "zh",
			deepl:  "ZH",
			libre:  "zh",
		},
		{
			name:   "swedish",
			lang:   "sv",
			google: "sv",
			deepl:  "SV",
			libre:  "sv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapLanguageToGoogle(tt.lang); got != tt.google {
				t.Errorf("mapLanguageToGoogle(%v) = %v, want %v", tt.lang, got, tt.google)
			}
			if got := mapLanguageToDeepL(tt.lang); got != tt.deepl {
				t.Errorf("mapLanguageToDeepL(%v) = %v, want %v", tt.lang, got, tt.deepl)
			}
			if got := mapLanguageToLibre(tt.lang); got != tt.libre {
				t.Errorf("mapLanguageToLibre(%v) = %v, want %v", tt.lang, got, tt.libre)
			}
		})
	}
}
