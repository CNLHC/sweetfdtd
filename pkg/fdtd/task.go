package fdtd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	logging "github.com/sweetfdtd/pkg/log"
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

var logger = logging.GetLogger()

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
	Wg          *sync.WaitGroup
	onFDTDLine  []func(line FDTDLine)
}

func (c *FDTDSlurmTask) OnFDTDLine(cb func(line FDTDLine)) {
	c.onFDTDLine = append(c.onFDTDLine, cb)
}

func (c *FDTDSlurmTask) Submit() {
	type Req map[string]interface{}
	cwd, _ := os.Getwd()
	payload, _ := json.Marshal(Req{
		"Script":      "#!/bin/bash\n /root/test_util",
		"UserId":      os.Getuid(),
		"GroupId":     os.Getegid(),
		"StdOut":      "test.out",
		"Name":        "ptest",
		"WorkDir":     cwd,
		"Environment": os.Environ(),
		"EnvSize":     len(os.Environ()),
	})
	c.SubmitTime = time.Now()
	c.SubmitCount++
	if res, err := slurm.SubmitBatchJob(payload); err != nil {
		c.Status.Slurm = SlurmJobExited
		logger.Errorf("can not submit job %v", err)
	} else {
		c.Status.Slurm = SlurmJobRunning
		if r, ok := (*res)["JobId"]; ok {
			c.JOBID = uint32(r.(uint))
			logger.Infof("submit JOB, ID: %d", r)
			go c.updateStatus()
		} else {
			logger.WithFields(logrus.Fields{
				"fsp":   c.FSPFile,
				"JobID": c.JOBID,
				"res":   fmt.Sprintf("%+v", *res),
			}).Errorf("No job id returned")

		}
	}
}
func waitfile(fp string) error {
	tick := time.NewTicker(time.Millisecond * 200)
	var max_retries = 20
	for _ = range tick.C {
		if _, err := os.Stat(fp); err != nil {
			if max_retries <= 0 {
				return errors.New("Wait for stdout file timeout")
			}
			if os.IsNotExist(err) {
				max_retries--
				continue
			} else {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

func (c *FDTDSlurmTask) updateStatus() {
	defer c.Wg.Done()
	local_fields := logrus.Fields{
		"fsp":   c.FSPFile,
		"JobID": c.JOBID,
	}

	if res, err := c.getJobInfo(); err != nil {
		logger.WithFields(local_fields).Errorf("can not find slurm job ")
		c.Status.Slurm = SlurmJobExited
	} else {
		logger.WithFields(local_fields).Infof("start monitoring")
		c.Status.Slurm = "Normal"
		cwd := res.JobArray[0].WorkDir
		stdout := res.JobArray[0].StdOut
		fp := path.Join(cwd+"/", stdout)
		if err := waitfile(fp); err != nil {
			logger.
				WithError(err).
				WithFields(local_fields).
				Error("can not open stdout file")
		}
		if file, err := os.Open(fp); err != nil {
			logger.WithFields(logrus.Fields{
				"fsp":   c.FSPFile,
				"JobID": c.JOBID,
				"fp":    fp,
				"err":   err.Error(),
			}).Errorf("error open stdout file")
			c.Status.FDTD = FDTDFailed
		} else {
			defer file.Close()
			ticker := time.NewTicker(200 * time.Millisecond)
			reader := bufio.NewReader(file)
		loop:
			for range ticker.C {
				if line, err := reader.ReadString('\n'); err != nil {
					if err != io.EOF {
						logger.WithFields(logrus.Fields{
							"fsp":   c.FSPFile,
							"JobID": c.JOBID,
							"err":   err.Error(),
						}).Error("read error")
					}
				} else {
					fdtdline := ParseStdOutLine(line)
					for _, fn := range c.onFDTDLine {
						if fn != nil {
							fn(fdtdline)
						}
					}
					switch fdtdline.Type {
					case LineTypeComplete:
						c.Status.FDTD = FDTDComplete
						break loop
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
