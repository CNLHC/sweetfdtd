package fdtd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

const (
	TaskUnschedule = iota
	TaskFailed     = iota
	TaskRunning    = iota
)

type FDTDSlurmTask struct {
	FSPFile      string
	JOBID        uint
	Status       int
	Progress     float32
	AutoShutfoff float32
	ETA          string
}

func (c *FDTDSlurmTask) Submit() string {
	cmd := exec.Command("sbatch")
	fmt.Printf(cmd.String())
	return cmd.String()
}

func (c *FDTDSlurmTask) UpdateStatus() bool {
	return false
}

func (c *FDTDSlurmTask) readLogFile() ([]byte, error) {
	return ioutil.ReadFile("")
}
