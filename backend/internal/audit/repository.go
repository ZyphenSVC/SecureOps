package audit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Record(ctx context.Context, actorUserID *string, action, resourceType string, resourceID *string, metadata string) error {
	query := `
		INSERT INTO audit_logs (
			actor_user_id,
			action,
			resource_type,
			resource_id,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`

	if metadata == "" {
		metadata = "{}"
	}

	_, err := r.db.Exec(ctx, query, actorUserID, action, resourceType, resourceID, metadata)
	return err
}

type LogEntry struct {
	ID           string
	ActorUserID  *string
	Action       string
	ResourceType string
	ResourceID   *string
	Metadata     string
	CreatedAt    string
}

func (r *Repository) List(ctx context.Context) ([]LogEntry, error) {
	query := `
		SELECT
			id,
			actor_user_id,
			action,
			resource_type,
			resource_id,
			metadata::text,
			created_at::text
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry

	for rows.Next() {
		var entry LogEntry

		err := rows.Scan(
			&entry.ID,
			&entry.ActorUserID,
			&entry.Action,
			&entry.ResourceType,
			&entry.ResourceID,
			&entry.Metadata,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
