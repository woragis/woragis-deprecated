package jobapplications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	jobApplicationQueueKey = "job-applications:queue"
	jobApplicationJobPrefix = "job-applications:job:"
)

// Queue manages job application jobs in Redis.
type Queue interface {
	EnqueueJob(ctx context.Context, job *JobApplicationJob) error
	DequeueJob(ctx context.Context, timeout time.Duration) (*JobApplicationJob, error)
	GetJob(ctx context.Context, jobID string) (*JobApplicationJob, error)
	MarkJobComplete(ctx context.Context, jobID string) error
	MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error
}

type redisQueue struct {
	client *redis.Client
}

// NewRedisQueue creates a new Redis-backed job application queue.
func NewRedisQueue(client *redis.Client) Queue {
	return &redisQueue{client: client}
}

func (q *redisQueue) EnqueueJob(ctx context.Context, job *JobApplicationJob) error {
	if q.client == nil {
		return NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Store job data
	jobKey := jobApplicationJobPrefix + job.ID
	if err := q.client.Set(ctx, jobKey, jobData, 24*time.Hour).Err(); err != nil {
		return NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	// Add to queue
	if err := q.client.LPush(ctx, jobApplicationQueueKey, job.ID).Err(); err != nil {
		return NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	return nil
}

func (q *redisQueue) DequeueJob(ctx context.Context, timeout time.Duration) (*JobApplicationJob, error) {
	if q.client == nil {
		return nil, NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	// Blocking pop from queue
	result, err := q.client.BRPop(ctx, timeout, jobApplicationQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No job available
		}
		return nil, NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	if len(result) < 2 {
		return nil, nil
	}

	jobID := result[1]
	return q.GetJob(ctx, jobID)
}

func (q *redisQueue) GetJob(ctx context.Context, jobID string) (*JobApplicationJob, error) {
	if q.client == nil {
		return nil, NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	jobKey := jobApplicationJobPrefix + jobID
	jobData, err := q.client.Get(ctx, jobKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, NewDomainError(ErrCodeNotFound, "jobapplications: job not found")
		}
		return nil, NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	var job JobApplicationJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

func (q *redisQueue) MarkJobComplete(ctx context.Context, jobID string) error {
	if q.client == nil {
		return NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
	}

	// Remove job from storage
	jobKey := jobApplicationJobPrefix + jobID
	return q.client.Del(ctx, jobKey).Err()
}

func (q *redisQueue) MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error {
	// For now, just log the failure. Could store in a failed jobs list.
	// The job will be removed from storage after a timeout
	return nil
}

