package api

import (
	"github.com/Aitugan/CodingChallenge/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func New() (*gin.Engine, error) {
	r := gin.Default()
	wm := handlers.NewWorksManager()
	work := r.Group("/work")
	{
		work.POST("/start", wm.Start)
		work.GET("/log/:id", wm.Log)
		work.PUT("/stop/:id", wm.Stop)
		work.GET("/query-status/:id", wm.QueryStatus)
	}

	return r, nil
}
