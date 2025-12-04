package application

import (
	"context"

	dto "github.com/dementievme/pull-request-service/internal/application/dto"
	entity "github.com/dementievme/pull-request-service/internal/domain/entity"
	repo "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type UserUseCase struct {
	repo repo.UserRepository
}

func NewUserUseCase(userRepo repo.UserRepository) *UserUseCase {
	return &UserUseCase{repo: userRepo}
}

func (u *UserUseCase) SetActive(ctx context.Context, userDTO *dto.SetActiveUserDTO) (*dto.UserDTO, error) {
	if err := u.repo.UpdateIsActive(ctx, userDTO.UserID, userDTO.IsActive); err != nil {
		return nil, err
	}

	userEntity, err := u.repo.GetByID(ctx, userDTO.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		UserID:   userEntity.ID,
		Username: userEntity.Name,
		TeamName: userEntity.TeamName,
		IsActive: userEntity.IsActive,
	}, nil
}

func (u *UserUseCase) GetUserByID(ctx context.Context, userID string) (*dto.UserDTO, error) {
	userEntity, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		UserID:   userEntity.ID,
		Username: userEntity.Name,
		TeamName: userEntity.TeamName,
		IsActive: userEntity.IsActive,
	}, nil
}

func (u *UserUseCase) Create(ctx context.Context, teamDTO *dto.TeamDTO) error {
	users := make([]*entity.User, 0, len(teamDTO.TeamMembers))

	for _, member := range teamDTO.TeamMembers {
		users = append(users, &entity.User{
			ID:       member.UserID,
			Name:     member.UserName,
			TeamName: teamDTO.TeamName,
			IsActive: member.IsActive,
		})
	}

	return u.repo.Create(ctx, users)
}
