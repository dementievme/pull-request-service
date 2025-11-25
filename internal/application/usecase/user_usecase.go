package application

import (
	"context"

	application "github.com/dementievme/pull-request-service/internal/application/dto"
	domain "github.com/dementievme/pull-request-service/internal/domain/repository"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
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
