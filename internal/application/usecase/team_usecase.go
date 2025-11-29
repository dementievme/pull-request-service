package application

import (
	"context"

	application "github.com/dementievme/pull-request-service/internal/application/dto"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type TeamUseCase struct {
	repo domain.TeamRepository
}

func NewTeamUseCase(repo domain.TeamRepository) *TeamUseCase {
	return &TeamUseCase{repo: repo}
}

func (u *TeamUseCase) CreateTeam(ctx context.Context, dto *application.TeamDTO) (*application.TeamDTO, error) {
	team := &entity.Team{
		Name: dto.TeamName,
	}
	if err := u.repo.Create(ctx, team); err != nil {
		return nil, err
	}
	return dto, nil
}

func (u *TeamUseCase) GetTeam(ctx context.Context, name string) (*application.TeamDTO, error) {
	t, err := u.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	dto := &application.TeamDTO{
		TeamName: t.Name,
	}
	for _, m := range t.Members {
		dto.TeamMembers = append(dto.TeamMembers, &application.TeamMemberDTO{
			UserID:   m.ID,
			UserName: m.Name,
			IsActive: m.IsActive,
		})
	}
	return dto, nil
}
