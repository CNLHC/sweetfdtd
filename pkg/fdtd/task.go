package fdtd

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/prometheus/common/log"
	"github.com/sweetfdtd/pkg/slurm"
)

const (
	FDTDFailed   = iota
	FDTDRunning  = iota
	FDTDComplete = iota
)
const (
	SlurmJobExited  = "Exit"
	SlurmJobRunning = "Running"
)

type FDTDSlurmTask struct {
	FSPFile string
	JOBID   uint32
	Status  struct {
		Slurm     string
		FDTD      int
		Calculate string
	}
	SubmitCount uint32
	SubmitTime  time.Time
	Statistic   LineProgress
}

func (c *FDTDSlurmTask) Submit() {
	type Res struct {
		Scripts string
	}
	payload, _ := json.Marshal(Res{
		Scripts: "#!/bin/bash\n#SBATCH -o stest%j.out\n hostname",
	})

	c.SubmitTime = time.Now()
	c.SubmitCount++
	if res, err := slurm.SubmitBatchJob(payload); err != nil {
		c.Status.Slurm = SlurmJobExited
		log.Errorf("can not submit job")
	} else {
		c.Status.Slurm = SlurmJobRunning
		if r, ok := (*res)["JobId"]; ok {
			c.JOBID = r.(uint32)
			log.Errorf("submit JOB, ID: %d", r)
		}
		go c.updateStatus()
	}
}

func (c *FDTDSlurmTask) updateStatus() {
	if res, err := c.getJobInfo(); err != nil {
		log.Warnf("can not find slurm job %d", c.JOBID)
		c.Status.Slurm = SlurmJobExited
	} else {
		c.Status.Slurm = "Normal"
		if file, err := os.Open(res.JobArray[0].StdOut); err != nil {
			c.Status.FDTD = FDTDFailed
		} else {
			defer file.Close()
			reader := bufio.NewReader(file)
			for {
				if line, err := reader.ReadString('\n'); err != nil {
					break
				} else {
					fdtdline := ParseStdOutLine(line)
					switch fdtdline.Type {
					case LineTypeComplete:
						c.Status.FDTD = FDTDComplete
						break
					case LineTypePlain:
						if fdtdline.PlainError == ErrorLicense {
							c.Status.FDTD = FDTDFailed
						}
						break
					case LineTypeUpdate:
						c.Statistic = fdtdline.Update
						break
					case LineTypeStatus:
						c.Status.Calculate = fdtdline.Status
						break
					}
				}
			}
		}
	}
}

func (c *FDTDSlurmTask) getJobInfo() (slurm.JobInfoMsg, error) {
	return slurm.LoadJobs(slurm.LoadJobsPayload{
		JobId: &c.JOBID,
	})
}
