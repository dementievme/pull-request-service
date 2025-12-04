package api

import (
	"net/http"

	dto "github.com/dementievme/pull-request-service/internal/application/dto"
	errors "github.com/dementievme/pull-request-service/internal/application/errors"
	application "github.com/dementievme/pull-request-service/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

func createTeamHandler(teamUseCase *application.TeamUseCase, userUseCase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var team dto.TeamDTO
		if err := c.ShouldBindJSON(&team); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "invalid request body")
			return
		}

		if _, err := teamUseCase.CreateTeam(c, &team); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrTeamExists.Error(), "team_name already exists")
			return
		}

		if err := userUseCase.Create(c, &team); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "failed to create team members")
			return
		}

		c.JSON(http.StatusCreated, gin.H{"team": team})
	}
}

func getTeamHandler(useCase *application.TeamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("team_name")
		team, err := useCase.GetTeam(c, name)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "team not found")
			return
		}
		c.JSON(http.StatusOK, team)
	}
}

func setUserActiveHandler(useCase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user dto.SetActiveUserDTO
		if err := c.ShouldBindJSON(&user); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "invalid request body")
			return
		}

		updated, err := useCase.SetActive(c, &user)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "user not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": updated})
	}
}

func createPRHandler(prUseCase *application.PullRequestUseCase, userUseCase *application.UserUseCase, teamUseCase *application.TeamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pr dto.CreatePullRequestDTO
		if err := c.ShouldBindJSON(&pr); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "invalid request body")
			return
		}

		user, err := userUseCase.GetUserByID(c, pr.AuthorID)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "author not found")
			return
		}

		teamDTO, err := teamUseCase.FindActiveUsers(c, user.TeamName)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "team not found")
			return
		}

		prCreated, err := prUseCase.CreatePR(c, &pr, teamDTO)
		if err != nil {
			sendError(c, http.StatusConflict, errors.ErrPRExists.Error(), "PR id already exists")
			return
		}

		c.JSON(http.StatusCreated, gin.H{"pr": prCreated})
	}
}

func mergePRHandler(useCase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var merge dto.MergePullRequestDTO
		if err := c.ShouldBindJSON(&merge); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "invalid request body")
			return
		}

		pr, err := useCase.Merge(c, merge.PullRequestID)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "PR not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{"pr": pr})
	}
}

func reassignPRHandler(prUseCase *application.PullRequestUseCase, userUseCase *application.UserUseCase, teamUseCase *application.TeamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var prDTO dto.ReassignPullRequestDTO
		if err := c.ShouldBindJSON(&prDTO); err != nil {
			sendError(c, http.StatusBadRequest, errors.ErrNotFound.Error(), "invalid request body")
			return
		}

		user, err := userUseCase.GetUserByID(c, prDTO.OldReviewerID)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "user not found")
			return
		}

		team, err := teamUseCase.FindActiveUsers(c, user.TeamName)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "team not found")
			return
		}

		pr, replacedBy, err := prUseCase.Reassign(c, &prDTO, team.TeamMembers)
		if err != nil {
			switch err.Error() {
			case errors.ErrPRMerged.Error():
				sendError(c, http.StatusConflict, errors.ErrPRMerged.Error(), "cannot reassign on merged PR")
			case errors.ErrNotAssigned.Error():
				sendError(c, http.StatusConflict, errors.ErrNotAssigned.Error(), "reviewer is not assigned to this PR")
			case errors.ErrNoCandidate.Error():
				sendError(c, http.StatusConflict, errors.ErrNoCandidate.Error(), "no active replacement candidate in team")
			default:
				sendError(c, http.StatusConflict, "UNKNOWN_ERROR", err.Error())
			}
			return
		}

		c.JSON(http.StatusOK, dto.ReassignPullRequestResponseDTO{
			PR:         *pr,
			ReplacedBy: replacedBy,
		})
	}
}

func getUserReviewsHandler(useCase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		prs, err := useCase.GetForReviewer(c, userID)
		if err != nil {
			sendError(c, http.StatusNotFound, errors.ErrNotFound.Error(), "user not found")
			return
		}

		c.JSON(http.StatusOK, dto.UserPRsDTO{
			UserID:       userID,
			PullRequests: prs,
		})
	}
}
