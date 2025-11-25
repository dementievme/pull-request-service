package infrastructure

import (
	"context"
	"database/sql"

	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
)

type PostgreTeamRepository struct {
	db *sql.DB
}

func NewPostgreTeamRepository(db *sql.DB) *PostgreTeamRepository {
	return &PostgreTeamRepository{db: db}
}

func (r *PostgreTeamRepository) Create(ctx context.Context, team *entity.Team) error {
	for _, m := range team.Members {
		_, err := r.db.ExecContext(ctx,
			`INSERT INTO users (id, name, is_active, team_name) 
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, is_active=EXCLUDED.is_active, team_name=EXCLUDED.team_name`,
			m.ID, m.Name, m.IsActive, team.Name,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgreTeamRepository) GetByName(ctx context.Context, name string) (*entity.Team, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, is_active FROM users WHERE team_name=$1`, name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*entity.User
	for rows.Next() {
		u := &entity.User{TeamName: name}
		if err := rows.Scan(&u.ID, &u.Name, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	if len(members) == 0 {
		return nil, ErrNotFound
	}

	return &entity.Team{
		Name:    name,
		Members: members,
	}, nil
}

func (r *PostgreTeamRepository) FindActiveUsers(ctx context.Context, teamName string) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name FROM users WHERE team_name=$1 AND is_active=true`, teamName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{TeamName: teamName, IsActive: true}
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
