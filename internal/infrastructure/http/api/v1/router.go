package api

import (
	usecase "github.com/dementievme/pull-request-service/internal/application/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	usecases *usecase.UseCases,
) {
	r.Use(cors.Default())

	teams := r.Group("/team")
	{
		teams.POST("/add", createTeamHandler(usecases.TeamUseCase, usecases.UserUseCase))
		teams.GET("/get", getTeamHandler(usecases.TeamUseCase))
	}

	users := r.Group("/users")
	{
		users.POST("/setIsActive", setUserActiveHandler(usecases.UserUseCase))
		users.GET("/getReview", getUserReviewsHandler(usecases.PullRequestUseCase))
	}

	prs := r.Group("/pullRequest")
	{
		prs.POST("/create", createPRHandler(usecases.PullRequestUseCase, usecases.UserUseCase, usecases.TeamUseCase))
		prs.POST("/merge", mergePRHandler(usecases.PullRequestUseCase))
		prs.POST("/reassign", reassignPRHandler(usecases.PullRequestUseCase, usecases.UserUseCase, usecases.TeamUseCase))
	}
}
