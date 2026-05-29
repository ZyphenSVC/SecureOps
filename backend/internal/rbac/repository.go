package rbac

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

func (r *Repository) UserHasPermission(ctx context.Context, userID string, permissionKey string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN role_permissions rp ON rp.role_id = ur.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE ur.user_id = $1
			  AND p.key = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, permissionKey).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
