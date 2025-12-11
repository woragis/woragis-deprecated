package resumes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	resumeQueueKey        = "resumes:queue"
	resumeJobPrefix       = "resumes:job:"
	resumeDeadLetterKey   = "resumes:dead-letter:queue"
	resumeEventsChannel   = "resumes:events"
)

// Queue manages resume generation jobs in Redis.
type Queue interface {
	EnqueueJob(ctx context.Context, job *ResumeJob) error
	GetJob(ctx context.Context, jobID string) (*ResumeJob, error)
	UpdateJobStatus(ctx context.Context, jobID string, status string, errorMsg *string, errorType *string, retryCount *int, result *ResumeJobResult) error
}

// ResumeJob represents a resume generation job.
type ResumeJob struct {
	ID              string    `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	JobApplicationID uuid.UUID `json:"job_application_id"`
	JobDescription  string    `json:"job_description"`
	JobTitle        string    `json:"job_title"`
	Language        string    `json:"language"`
	Status          string    `json:"status"` // pending, processing, completed, failed, retrying, dead_letter
	RetryCount      int       `json:"retry_count"`
	MaxRetries      int       `json:"max_retries"`
	LastError       string    `json:"last_error,omitempty"`
	LastErrorType   string    `json:"last_error_type,omitempty"`
	LastErrorAt     time.Time `json:"last_error_at,omitempty"`
	Result          *ResumeJobResult `json:"result,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ResumeJobResult contains the result of a completed resume generation job.
type ResumeJobResult struct {
	OutputPath string   `json:"output_path"`
	FileName   string   `json:"file_name"`
	FileSize   int64    `json:"file_size"`
	Tags       []string `json:"tags,omitempty"`
	DurationMs int     `json:"duration_ms,omitempty"`
}

type redisQueue struct {
	client *redis.Client
}

// NewRedisQueue creates a new Redis-backed resume generation queue.
func NewRedisQueue(client *redis.Client) Queue {
	return &redisQueue{client: client}
}

func (q *redisQueue) EnqueueJob(ctx context.Context, job *ResumeJob) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
	}

	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	job.Status = "pending"
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	if job.MaxRetries == 0 {
		job.MaxRetries = 3
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Store job data with 7 day TTL
	jobKey := resumeJobPrefix + job.ID
	if err := q.client.Set(ctx, jobKey, jobData, 7*24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to store job: %w", err)
	}

	// Add to queue
	if err := q.client.LPush(ctx, resumeQueueKey, job.ID).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

func (q *redisQueue) GetJob(ctx context.Context, jobID string) (*ResumeJob, error) {
	if q.client == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	jobKey := resumeJobPrefix + jobID
	jobData, err := q.client.Get(ctx, jobKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	var job ResumeJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

func (q *redisQueue) UpdateJobStatus(ctx context.Context, jobID string, status string, errorMsg *string, errorType *string, retryCount *int, result *ResumeJobResult) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
	}

	job, err := q.GetJob(ctx, jobID)
	if err != nil {
		return err
	}

	job.Status = status
	job.UpdatedAt = time.Now()

	if errorMsg != nil {
		job.LastError = *errorMsg
		job.LastErrorAt = time.Now()
	}

	if errorType != nil {
		job.LastErrorType = *errorType
	}

	if retryCount != nil {
		job.RetryCount = *retryCount
	}

	if result != nil {
		job.Result = result
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	jobKey := resumeJobPrefix + jobID
	if err := q.client.Set(ctx, jobKey, jobData, 7*24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	// Publish event
	event := map[string]interface{}{
		"job_id":  jobID,
		"status":  status,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if errorMsg != nil {
		event["error"] = *errorMsg
	}
	if result != nil {
		event["result"] = result
	}

	eventData, _ := json.Marshal(event)
	q.client.Publish(ctx, resumeEventsChannel, eventData)

	return nil
}

