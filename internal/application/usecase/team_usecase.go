package application

import (
	"context"

	dto "github.com/dementievme/pull-request-service/internal/application/dto"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type TeamUseCase struct {
	repo domain.TeamRepository
}

func NewTeamUseCase(teamRepo domain.TeamRepository) *TeamUseCase {
	return &TeamUseCase{repo: teamRepo}
}

func (u *TeamUseCase) CreateTeam(ctx context.Context, teamReq *dto.TeamDTO) (*dto.TeamDTO, error) {
	team := &entity.Team{
		Name: teamReq.TeamName,
	}

	if err := u.repo.Create(ctx, team); err != nil {
		return nil, err
	}

	return teamReq, nil
}

func (u *TeamUseCase) GetTeam(ctx context.Context, teamName string) (*dto.TeamDTO, error) {
	teamEntity, err := u.repo.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	teamDTO := &dto.TeamDTO{
		TeamName: teamEntity.Name,
	}

	for _, member := range teamEntity.Members {
		teamDTO.TeamMembers = append(teamDTO.TeamMembers, &dto.TeamMemberDTO{
			UserID:   member.ID,
			UserName: member.Name,
			IsActive: member.IsActive,
		})
	}

	return teamDTO, nil
}

func (u *TeamUseCase) FindActiveUsers(ctx context.Context, teamName string) (*dto.TeamDTO, error) {
	activeUsers, err := u.repo.FindActiveUsers(ctx, teamName)
	if err != nil {
		return nil, err
	}

	result := &dto.TeamDTO{
		TeamName: teamName,
	}

	for _, user := range activeUsers {
		result.TeamMembers = append(result.TeamMembers, &dto.TeamMemberDTO{
			UserID:   user.ID,
			UserName: user.Name,
			IsActive: user.IsActive,
		})
	}

	return result, nil
}
