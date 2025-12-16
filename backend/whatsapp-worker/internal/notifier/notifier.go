package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3" // SQLite driver for whatsmeow
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Notifier defines the contract for sending WhatsApp messages.
type Notifier interface {
	Send(ctx context.Context, to string, message string) error
	Connect(ctx context.Context) error
	Disconnect()
	IsConnected() bool
	GetQRCode() string
}

// WhatsmeowNotifier implements WhatsApp notifications using whatsmeow.
type WhatsmeowNotifier struct {
	client      *whatsmeow.Client
	logger      *slog.Logger
	mu          sync.RWMutex
	connected   bool
	sessionPath string
	qrCode      string
	qrCodeMu    sync.RWMutex
}

// NewWhatsmeowNotifier creates a new whatsmeow-based notifier.
// sessionPath is the directory where the WhatsApp session will be stored.
func NewWhatsmeowNotifier(sessionPath string, logger *slog.Logger) (*WhatsmeowNotifier, error) {
	if sessionPath == "" {
		sessionPath = "./whatsapp-session"
	}

	// Ensure session directory exists
	if err := os.MkdirAll(sessionPath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	dbPath := filepath.Join(sessionPath, "whatsapp.db")
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	container, err := sqlstore.New(context.Background(), "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath), dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	notifier := &WhatsmeowNotifier{
		client:      client,
		logger:      logger,
		sessionPath: sessionPath,
		connected:   false,
	}

	// Set up event handlers
	client.AddEventHandler(notifier.eventHandler)

	return notifier, nil
}

// Connect connects to WhatsApp. If not logged in, it will generate a QR code.
// The QR code will be logged to the logger.
func (n *WhatsmeowNotifier) Connect(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.connected {
		return nil
	}

	if n.client.Store.ID == nil {
		// No ID stored, need to login with QR code
		if n.logger != nil {
			n.logger.Info("whatsapp: no session found, generating QR code for login")
		}

		qrChan, _ := n.client.GetQRChannel(ctx)
		err := n.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		// Wait for QR code in a goroutine
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					// Store QR code for API access
					n.qrCodeMu.Lock()
					n.qrCode = evt.Code
					n.qrCodeMu.Unlock()

					if n.logger != nil {
						n.logger.Info("whatsapp: QR code generated",
							slog.String("qr_code", evt.Code),
							slog.String("hint", "Scan this QR code with WhatsApp on your phone"))
					}
					// Also print to stdout for easy access
					fmt.Printf("\n=== WhatsApp QR Code ===\n%s\n========================\n", evt.Code)
				} else if evt.Event == "success" {
					// Clear QR code on success
					n.qrCodeMu.Lock()
					n.qrCode = ""
					n.qrCodeMu.Unlock()

					n.mu.Lock()
					n.connected = true
					n.mu.Unlock()
					if n.logger != nil {
						n.logger.Info("whatsapp: successfully logged in")
					}
				} else if evt.Event == "timeout" {
					// Clear QR code on timeout
					n.qrCodeMu.Lock()
					n.qrCode = ""
					n.qrCodeMu.Unlock()

					if n.logger != nil {
						n.logger.Warn("whatsapp: QR code timeout, please restart to generate new QR code")
					}
				}
			}
		}()
	} else {
		// Already logged in, just connect
		err := n.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		n.connected = true
		if n.logger != nil {
			n.logger.Info("whatsapp: connected with existing session")
		}
	}

	return nil
}

// Disconnect disconnects from WhatsApp.
func (n *WhatsmeowNotifier) Disconnect() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.client != nil {
		n.client.Disconnect()
		n.connected = false
		if n.logger != nil {
			n.logger.Info("whatsapp: disconnected")
		}
	}
}

// Send sends a WhatsApp message to the specified phone number.
// The phone number should be in international format (e.g., "5511999999999").
func (n *WhatsmeowNotifier) Send(ctx context.Context, to string, message string) error {
	n.mu.RLock()
	connected := n.connected
	n.mu.RUnlock()

	if !connected {
		// Try to connect if not connected
		// Note: Connection may require QR code scan, so we return an error
		// The message will be requeued and can be retried later
		return fmt.Errorf("whatsapp not connected - please ensure WhatsApp is authenticated via QR code")
	}

	// Parse phone number (remove any non-digit characters except +)
	var phoneNumber types.JID
	if len(to) > 0 && to[0] == '+' {
		phoneNumber = types.NewJID(to[1:], types.DefaultUserServer)
	} else {
		phoneNumber = types.NewJID(to, types.DefaultUserServer)
	}

	// Send the message - construct a text message
	msg := &waE2E.Message{
		Conversation: &message,
	}
	_, err := n.client.SendMessage(ctx, phoneNumber, msg)
	if err != nil {
		if n.logger != nil {
			n.logger.Error("whatsapp: failed to send message",
				slog.String("to", to),
				slog.String("error", err.Error()))
		}
		return fmt.Errorf("failed to send message: %w", err)
	}

	if n.logger != nil {
		n.logger.Info("whatsapp: message sent",
			slog.String("to", to),
			slog.String("message_preview", truncateString(message, 50)))
	}

	return nil
}

// eventHandler handles WhatsApp events.
func (n *WhatsmeowNotifier) eventHandler(evt interface{}) {
	switch evt.(type) {
	case *events.Connected:
		n.mu.Lock()
		n.connected = true
		n.mu.Unlock()
		if n.logger != nil {
			n.logger.Info("whatsapp: connected")
		}
	case *events.Disconnected:
		n.mu.Lock()
		n.connected = false
		n.mu.Unlock()
		if n.logger != nil {
			n.logger.Warn("whatsapp: disconnected")
		}
	case *events.LoggedOut:
		n.mu.Lock()
		n.connected = false
		n.mu.Unlock()
		if n.logger != nil {
			n.logger.Warn("whatsapp: logged out")
		}
	}
}

// GetQRCode returns the current QR code if available.
func (n *WhatsmeowNotifier) GetQRCode() string {
	n.qrCodeMu.RLock()
	defer n.qrCodeMu.RUnlock()
	return n.qrCode
}

// IsConnected returns whether the notifier is connected to WhatsApp.
func (n *WhatsmeowNotifier) IsConnected() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.connected
}

// truncateString truncates a string to a maximum length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
