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
		if res, err := LoadJobs(payload); err != nil {
			ctx.JSON(500, gin.H{"errMsg": err.Error()})
		} else {
			ctx.JSON(200, res)
		}
	})

	engine.GET("/job/submit", func(ctx *gin.Context) {
		if raw, err := ctx.GetRawData(); err != nil {
			ctx.JSON(500, gin.H{"errMsg": err.Error()})
		} else {
			if res, err := SubmitBatchJob(raw); err != nil {
				ctx.JSON(500, gin.H{"errMsg": err.Error()})
			} else {
				ctx.JSON(200, res)
			}
		}
	})
	engine.Run(":8443")
}
