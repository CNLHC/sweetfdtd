package slurm

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Listen() {
	engine := gin.Default()
	engine.GET("/jobs", func(ctx *gin.Context) {
		var payload LoadJobsPayload
		ctx.MustBindWith(&payload, binding.JSON)
		res := LoadJobs(payload)
		ctx.JSON(200, res)
	})
	engine.Run(":8443")
}
