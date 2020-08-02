package fdtd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sweetfdtd/pkg/slurm"
)

const (
	TaskUnschedule = iota
	TaskFailed     = iota
	TaskRunning    = iota
)

type FDTDSlurmTask struct {
	FSPFile      string
	JOBID        uint
	UserId       uint32
	Status       int
	Progress     float32
	AutoShutfoff float32
	ETA          string
}

func (c *FDTDSlurmTask) Submit() {
	type Res struct {
		Scripts string
	}
	payload, _ := json.Marshal(Res{
		Scripts: "hi",
	})
	res, _ := slurm.SubmitBatchJob(payload)
	if r, ok := (*res)["JobId"]; ok {
		c.JOBID = r.(uint)
	}
}

func (c *FDTDSlurmTask) UpdateStatus() bool {
	return false
}

func (c *FDTDSlurmTask) readLogFile() ([]byte, error) {
	return ioutil.ReadFile("")
}
