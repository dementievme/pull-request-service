package application

import (
	"context"

	application "github.com/dementievme/pull-request-service/internal/application/dto"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	repo "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type UserUseCase struct {
	repo repo.UserRepository
}

func NewUserUseCase(repo repo.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) SetActive(ctx context.Context, dto *application.SetActiveUserDTO) (*application.UserDTO, error) {
	if err := u.repo.UpdateIsActive(ctx, dto.UserID, dto.IsActive); err != nil {
		return nil, err
	}
	user, err := u.repo.GetByID(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}
	return &application.UserDTO{
		UserID:   user.ID,
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}, nil
}

func (u *UserUseCase) Create(ctx context.Context, dto *application.TeamDTO) error {
	users := make([]*entity.User, 0, len(dto.TeamMembers))

	for _, member := range dto.TeamMembers {
		users = append(users, &entity.User{
			ID:       member.UserID,
			Name:     member.UserName,
			TeamName: dto.TeamName,
			IsActive: member.IsActive,
		})
	}

	err := u.repo.Create(ctx, users)
	if err != nil {
		return err
	}

	return nil
}
