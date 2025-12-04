package api

import (
	"net/http"

	dto "github.com/dementievme/pull-request-service/internal/application/dto"
	application "github.com/dementievme/pull-request-service/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

func createTeamHandler(teamUseCase *application.TeamUseCase, userUseCase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto dto.TeamDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := teamUseCase.CreateTeam(c, &dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userUseCase.Create(c, &dto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"team": dto})
	}
}

func getTeamHandler(usecase *application.TeamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("team_name")
		team, err := usecase.GetTeam(c, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, team)
	}
}

func setUserActiveHandler(usecase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto dto.SetActiveUserDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := usecase.SetActive(c, &dto)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func createPRHandler(usecase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto dto.CreatePullRequestDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pr, err := usecase.CreatePR(c, &dto)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"pr": pr})
	}
}

func mergePRHandler(usecase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto dto.MergePullRequestDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pr, err := usecase.Merge(c, dto.PullRequestID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pr": pr})
	}
}

func reassignPRHandler(usecase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ReassignPullRequestDTO
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pr, replacedBy, err := usecase.Reassign(c, &request)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.ReassignPullRequestResponseDTO{
			PR:         *pr,
			ReplacedBy: replacedBy,
		})
	}
}

func getUserReviewsHandler(usecase *application.PullRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		prs, err := usecase.GetForReviewer(c, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.UserPRsDTO{UserID: userID, PullRequests: prs})
	}
}
