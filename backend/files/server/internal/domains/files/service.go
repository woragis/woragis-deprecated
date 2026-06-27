package files

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// StorageProvider defines file storage operations.
type StorageProvider interface {
	Save(ctx context.Context, file *File, reader io.Reader) error
	Get(ctx context.Context, file *File) (io.ReadCloser, error)
	Delete(ctx context.Context, file *File) error
	GetURL(ctx context.Context, file *File) (string, error)
}

// Service orchestrates file workflows.
type Service interface {
	UploadFile(ctx context.Context, userID uuid.UUID, req UploadFileRequest) (*File, error)
	GetFile(ctx context.Context, fileID uuid.UUID) (*File, error)
	GetFileByUser(ctx context.Context, fileID, userID uuid.UUID) (*File, error)
	DeleteFile(ctx context.Context, fileID, userID uuid.UUID) error
	ListFiles(ctx context.Context, filters FileFilters) ([]File, error)
	DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, *File, error)
	GetFileURL(ctx context.Context, fileID uuid.UUID) (string, error)
}

// UploadFileRequest represents a file upload request.
type UploadFileRequest struct {
	FileName    string
	MimeType    string
	Size        int64
	Reader      io.Reader
	Metadata    map[string]interface{}
}

type service struct {
	repo            Repository
	storageProvider StorageProvider
	basePath        string
	logger          *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, storageProvider StorageProvider, basePath string, logger *slog.Logger) Service {
	return &service{
		repo:            repo,
		storageProvider: storageProvider,
		basePath:        basePath,
		logger:          logger,
	}
}

func (s *service) UploadFile(ctx context.Context, userID uuid.UUID, req UploadFileRequest) (*File, error) {
	// Generate file path
	fileID := uuid.New()
	ext := filepath.Ext(req.FileName)
	relativePath := filepath.Join(
		time.Now().Format("2006/01/02"),
		fileID.String()+ext,
	)
	fullPath := filepath.Join(s.basePath, relativePath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
	}

	// Create file entity
	file, err := NewFile(userID, req.FileName, relativePath, req.MimeType, req.Size)
	if err != nil {
		return nil, err
	}

	// Calculate hash while saving
	hash := sha256.New()
	multiWriter := io.MultiWriter(hash)

	// Save file using storage provider
	if s.storageProvider != nil {
		if err := s.storageProvider.Save(ctx, file, io.TeeReader(req.Reader, multiWriter)); err != nil {
			return nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
	} else {
		// Fallback to local file system
		outFile, err := os.Create(fullPath)
		if err != nil {
			return nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
		defer outFile.Close()

		if _, err := io.Copy(io.MultiWriter(outFile, multiWriter), req.Reader); err != nil {
			return nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
	}

	// Set hash
	file.Hash = hex.EncodeToString(hash.Sum(nil))

	// Check for duplicate
	existingFile, err := s.repo.GetFileByHash(ctx, file.Hash)
	if err == nil && existingFile != nil {
		// File already exists, return existing file
		return existingFile, nil
	}

	// Mark as ready
	file.MarkReady()

	// Save file metadata
	if err := s.repo.CreateFile(ctx, file); err != nil {
		return nil, err
	}

	// Get URL if storage provider supports it
	if s.storageProvider != nil {
		url, err := s.storageProvider.GetURL(ctx, file)
		if err == nil {
			file.URL = url
			_ = s.repo.UpdateFile(ctx, file)
		}
	}

	return file, nil
}

func (s *service) GetFile(ctx context.Context, fileID uuid.UUID) (*File, error) {
	return s.repo.GetFile(ctx, fileID)
}

func (s *service) GetFileByUser(ctx context.Context, fileID, userID uuid.UUID) (*File, error) {
	return s.repo.GetFileByUser(ctx, fileID, userID)
}

func (s *service) DeleteFile(ctx context.Context, fileID, userID uuid.UUID) error {
	file, err := s.repo.GetFileByUser(ctx, fileID, userID)
	if err != nil {
		return err
	}

	// Delete from storage
	if s.storageProvider != nil {
		if err := s.storageProvider.Delete(ctx, file); err != nil {
			s.logger.Warn("failed to delete file from storage", "error", err, "file_id", fileID)
		}
	} else {
		// Fallback to local file system
		fullPath := filepath.Join(s.basePath, file.Path)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			s.logger.Warn("failed to delete file from filesystem", "error", err, "file_id", fileID)
		}
	}

	// Delete from database
	return s.repo.DeleteFile(ctx, fileID, userID)
}

func (s *service) ListFiles(ctx context.Context, filters FileFilters) ([]File, error) {
	return s.repo.ListFiles(ctx, filters)
}

func (s *service) DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, *File, error) {
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return nil, nil, err
	}

	var reader io.ReadCloser
	if s.storageProvider != nil {
		reader, err = s.storageProvider.Get(ctx, file)
		if err != nil {
			return nil, nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
	} else {
		// Fallback to local file system
		fullPath := filepath.Join(s.basePath, file.Path)
		reader, err = os.Open(fullPath)
		if err != nil {
			return nil, nil, NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
	}

	return reader, file, nil
}

func (s *service) GetFileURL(ctx context.Context, fileID uuid.UUID) (string, error) {
	file, err := s.repo.GetFile(ctx, fileID)
	if err != nil {
		return "", err
	}

	if file.URL != "" {
		return file.URL, nil
	}

	if s.storageProvider != nil {
		url, err := s.storageProvider.GetURL(ctx, file)
		if err != nil {
			return "", NewDomainError(ErrCodeStorageFailure, ErrStorageError)
		}
		return url, nil
	}

	// Fallback: return relative path
	return "/files/" + file.ID.String(), nil
}

