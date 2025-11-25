package api

import (
	usecase "github.com/dementievme/pull-request-service/internal/application/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	pullRequestUseCase *usecase.PullRequestUseCase,
	teamUseCase *usecase.TeamUseCase,
	userUseCase *usecase.UserUseCase,
) {
	r.Use(cors.Default())

	teams := r.Group("/team")
	{
		teams.POST("/add", createTeamHandler(teamUseCase))
		teams.GET("/get", getTeamHandler(teamUseCase))
	}

	users := r.Group("/users")
	{
		users.POST("/setIsActive", setUserActiveHandler(userUseCase))
		users.GET("/getReview", getUserReviewsHandler(pullRequestUseCase))
	}

	prs := r.Group("/pullRequest")
	{
		prs.POST("/create", createPRHandler(pullRequestUseCase))
		prs.POST("/merge", mergePRHandler(pullRequestUseCase))
		prs.POST("/reassign", reassignPRHandler(pullRequestUseCase))
	}
}
