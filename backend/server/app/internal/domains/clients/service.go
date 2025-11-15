package clients

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service defines client domain operations consumed by transport layers.
type Service interface {
	CreateClient(ctx context.Context, req CreateClientRequest) (*Client, error)
	UpdateClient(ctx context.Context, req UpdateClientRequest) (*Client, error)
	GetClient(ctx context.Context, userID, id uuid.UUID) (*Client, error)
	ListClients(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]Client, error)
	DeleteClient(ctx context.Context, userID, id uuid.UUID) error
	ToggleArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error
	GetClientByPhoneNumber(ctx context.Context, userID uuid.UUID, phoneNumber string) (*Client, error)
}

// service orchestrates client domain use-cases.
type service struct {
	repo   Repository
	logger *slog.Logger
}

// Ensure service implements Service.
var _ Service = (*service)(nil)

// NewService builds a client domain service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// CreateClientRequest transports API payloads to the domain layer.
type CreateClientRequest struct {
	UserID      uuid.UUID
	Name        string
	Email       string
	PhoneNumber string
	Company     string
	Notes       string
}

// UpdateClientRequest handles partial updates.
type UpdateClientRequest struct {
	UserID      uuid.UUID
	ClientID    uuid.UUID
	Name        *string
	Email       *string
	PhoneNumber *string
	Company     *string
	Notes       *string
}

// CreateClient creates a new client.
func (s *service) CreateClient(ctx context.Context, req CreateClientRequest) (*Client, error) {
	client, err := NewClient(req.UserID, req.Name, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Company != "" {
		client.Company = req.Company
	}
	if req.Notes != "" {
		client.Notes = req.Notes
	}

	if err := client.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// UpdateClient updates an existing client.
func (s *service) UpdateClient(ctx context.Context, req UpdateClientRequest) (*Client, error) {
	client, err := s.repo.Get(ctx, req.UserID, req.ClientID)
	if err != nil {
		return nil, err
	}

	if err := client.UpdateMutableFields(req.Name, req.Email, req.PhoneNumber, req.Company, req.Notes); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// GetClient retrieves a client by ID.
func (s *service) GetClient(ctx context.Context, userID, id uuid.UUID) (*Client, error) {
	return s.repo.Get(ctx, userID, id)
}

// ListClients retrieves all clients for a user.
func (s *service) ListClients(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]Client, error) {
	return s.repo.List(ctx, userID, includeArchived)
}

// DeleteClient permanently removes a client.
func (s *service) DeleteClient(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.Delete(ctx, userID, id)
}

// ToggleArchived updates the archived status of a client.
func (s *service) ToggleArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error {
	return s.repo.SetArchived(ctx, userID, id, archived)
}

// GetClientByPhoneNumber finds a client by phone number.
func (s *service) GetClientByPhoneNumber(ctx context.Context, userID uuid.UUID, phoneNumber string) (*Client, error) {
	return s.repo.FindByPhoneNumber(ctx, userID, phoneNumber)
}

