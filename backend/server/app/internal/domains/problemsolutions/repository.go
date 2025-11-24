package problemsolutions

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for problem solutions.
type Repository interface {
	CreateProblemSolution(ctx context.Context, problemSolution *ProblemSolution) error
	UpdateProblemSolution(ctx context.Context, problemSolution *ProblemSolution) error
	GetProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) (*ProblemSolution, error)
	GetProblemSolutionPublic(ctx context.Context, problemSolutionID uuid.UUID) (*ProblemSolution, error)
	ListProblemSolutions(ctx context.Context, userID uuid.UUID) ([]ProblemSolution, error)
	ListFeaturedProblemSolutions(ctx context.Context) ([]ProblemSolution, error)
	DeleteProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) error
	GetProblemSolutionMatrix(ctx context.Context) ([]ProblemSolutionMatrixEntry, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateProblemSolution(ctx context.Context, problemSolution *ProblemSolution) error {
	if err := problemSolution.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(problemSolution).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateProblemSolution(ctx context.Context, problemSolution *ProblemSolution) error {
	if err := problemSolution.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(problemSolution).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) (*ProblemSolution, error) {
	var problemSolution ProblemSolution
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", problemSolutionID, userID).
		First(&problemSolution).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrProblemSolutionNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &problemSolution, nil
}

func (r *gormRepository) GetProblemSolutionPublic(ctx context.Context, problemSolutionID uuid.UUID) (*ProblemSolution, error) {
	var problemSolution ProblemSolution
	err := r.db.WithContext(ctx).
		Where("id = ?", problemSolutionID).
		First(&problemSolution).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrProblemSolutionNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &problemSolution, nil
}

func (r *gormRepository) ListProblemSolutions(ctx context.Context, userID uuid.UUID) ([]ProblemSolution, error) {
	var problemSolutions []ProblemSolution
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&problemSolutions).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return problemSolutions, nil
}

func (r *gormRepository) ListFeaturedProblemSolutions(ctx context.Context) ([]ProblemSolution, error) {
	var problemSolutions []ProblemSolution
	err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("created_at DESC").
		Find(&problemSolutions).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return problemSolutions, nil
}

func (r *gormRepository) DeleteProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) error {
	var problemSolution ProblemSolution
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", problemSolutionID, userID).
		First(&problemSolution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodeNotFound, ErrProblemSolutionNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&problemSolution).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

// ProblemSolutionMatrixEntry represents a technology and the problems solved with it.
type ProblemSolutionMatrixEntry struct {
	Technology string   `json:"technology"`
	Problems   []string `json:"problems"`
	Count      int      `json:"count"`
}

func (r *gormRepository) GetProblemSolutionMatrix(ctx context.Context) ([]ProblemSolutionMatrixEntry, error) {
	// Get all problem solutions
	var problemSolutions []ProblemSolution
	if err := r.db.WithContext(ctx).Find(&problemSolutions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Build technology map: technology -> []problem IDs
	techMap := make(map[string]map[string]bool) // technology -> problem ID set

	for _, ps := range problemSolutions {
		for _, tech := range ps.Technologies {
			if techMap[tech] == nil {
				techMap[tech] = make(map[string]bool)
			}
			techMap[tech][ps.ID.String()] = true
		}
	}

	// Convert to matrix entries
	var matrix []ProblemSolutionMatrixEntry
	for tech, problemIDs := range techMap {
		// Get problem summaries for this technology
		var problems []string
		for problemID := range problemIDs {
			// Find the problem solution to get the problem text
			for _, ps := range problemSolutions {
				if ps.ID.String() == problemID {
					// Use a short version of the problem (first 100 chars)
					problemText := ps.Problem
					if len(problemText) > 100 {
						problemText = problemText[:100] + "..."
					}
					problems = append(problems, problemText)
					break
				}
			}
		}

		matrix = append(matrix, ProblemSolutionMatrixEntry{
			Technology: tech,
			Problems:   problems,
			Count:      len(problems),
		})
	}

	return matrix, nil
}

