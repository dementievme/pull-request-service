package infrastructure

import (
	"context"
	"database/sql"

	entity "github.com/dementievme/pull-request-service/internal/domain/entity"

	_ "github.com/lib/pq"
)

type PostgreUserRepository struct {
	db *sql.DB
}

func NewPostgreUserRepository(db *sql.DB) *PostgreUserRepository {
	return &PostgreUserRepository{db: db}
}

func (r *PostgreUserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`, user.ID, user.Name, user.TeamName, user.IsActive)
	return err
}

func (r *PostgreUserRepository) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, team_name, is_active FROM users WHERE id = $1
	`, userID).Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgreUserRepository) UpdateIsActive(ctx context.Context, userID string, isActive bool) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE users SET is_active = $2 WHERE id = $1
	`, userID, isActive)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// FindByTeam возвращает всех пользователей команды
func (r *PostgreUserRepository) FindByTeam(ctx context.Context, teamName string) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, team_name, is_active FROM users WHERE team_name = $1 ORDER BY name
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if len(users) == 0 {
		return nil, ErrNotFound
	}
	return users, rows.Err()
}
