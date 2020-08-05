package fdtd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	logging "github.com/sweetfdtd/pkg/log"
	"github.com/sweetfdtd/pkg/slurm"
)

type TaskSlurmStatus int
type TaskFDTDStatus int

const (
	FDTDWaiting  TaskFDTDStatus = iota
	FDTDTimeout  TaskFDTDStatus = iota
	FDTDFailed   TaskFDTDStatus = iota
	FDTDRunning  TaskFDTDStatus = iota
	FDTDComplete TaskFDTDStatus = iota
)
const (
	SlurmJobExited  TaskSlurmStatus = iota
	SlurmJobRunning TaskSlurmStatus = iota
)

var logger = logging.GetLogger()

type FDTDSlurmTask struct {
	FSPFile string
	JOBID   uint32
	Status  struct {
		Slurm     TaskSlurmStatus
		FDTD      TaskFDTDStatus
		Calculate string
	}
	SubmitCount    uint32
	SubmitTime     time.Time
	Statistic      LineProgress
	Wg             *sync.WaitGroup
	onFDTDLineRecv []func(line FDTDLine)
	RetryInSec     int
}

func (c *FDTDSlurmTask) OnFDTDLineReceived(cb func(line FDTDLine)) {
	c.onFDTDLineRecv = append(c.onFDTDLineRecv, cb)
}

func (c *FDTDSlurmTask) Submit() {
	type Req map[string]interface{}
	cwd, _ := os.Getwd()
	payload, _ := json.Marshal(Req{
		"Script":      "#!/bin/bash\n /root/test_util",
		"UserId":      os.Getuid(),
		"GroupId":     os.Getegid(),
		"StdOut":      "test-%j.out",
		"Name":        "ptest",
		"WorkDir":     cwd,
		"ExcNodes":    "c2",
		"MaxCpus":     1,
		"Environment": os.Environ(),
		"EnvSize":     len(os.Environ()),
	})
	c.SubmitTime = time.Now()
	c.SubmitCount++
	if res, err := slurm.SubmitBatchJob(payload); err != nil {
		c.Status.Slurm = SlurmJobExited
		logger.Errorf("can not submit job %v", err)
		c.Wg.Done()
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
func (c *FDTDSlurmTask) retry() {
	ticker := time.NewTicker(time.Second)
	c.RetryInSec = 5
	for range ticker.C {
		c.RetryInSec--
		if c.RetryInSec <= 0 {
			break
		}
	}
	c.Submit()
}

func waitfile(fp string) error {
	logger.Infof("try to open fp:%s", fp)
	tick := time.NewTicker(time.Millisecond * 200)
	var max_retries = 2000
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

func (c *FDTDSlurmTask) setSlurmStatus(status TaskSlurmStatus) {
	c.Status.Slurm = status
}

func (c *FDTDSlurmTask) setFDTDStatus(status TaskFDTDStatus) {
	c.Status.FDTD = status

	if status == FDTDFailed || status == FDTDTimeout {
		c.Wg.Add(1)
		go c.retry()
	}
}

func (c *FDTDSlurmTask) updateStatus() {
	defer c.Wg.Done()
	local_fields := logrus.Fields{
		"fsp":   c.FSPFile,
		"JobID": c.JOBID,
	}
	var (
		res  slurm.JobInfoMsg
		file *os.File
		err  error
	)

	if res, err = c.getJobInfo(); err != nil {
		logger.WithFields(local_fields).Errorf("can not find slurm job ")
		c.setSlurmStatus(SlurmJobExited)
		return
	}

	c.setSlurmStatus(SlurmJobRunning)
	c.setFDTDStatus(FDTDWaiting)
	logger.WithFields(local_fields).Infof("start monitoring")

	//get stdout file
	cwd := res.JobArray[0].WorkDir
	stdout := res.JobArray[0].StdOut
	fp := path.Join(cwd+"/", stdout)
	fp = strings.Replace(fp, "%j", strconv.Itoa(int(c.JOBID)), 1)
	if err := waitfile(fp); err != nil {
		c.setFDTDStatus(FDTDTimeout)
		logger.
			WithError(err).
			WithFields(local_fields).
			Error("can not open stdout file")
		return
	}

	// try to open file
	if file, err = os.Open(fp); err != nil {
		logger.WithFields(logrus.Fields{
			"fsp":   c.FSPFile,
			"JobID": c.JOBID,
			"fp":    fp,
			"err":   err.Error(),
		}).Errorf("error open stdout file")
		c.setFDTDStatus(FDTDFailed)
		return
	}

	defer file.Close()
	c.setFDTDStatus(FDTDRunning)
	ticker := time.NewTicker(100 * time.Millisecond)
	reader := bufio.NewReader(file)
loop:
	for range ticker.C {
		var line string
		if line, err = reader.ReadString('\n'); err != nil {
			if err != io.EOF {
				logger.WithFields(logrus.Fields{
					"fsp":   c.FSPFile,
					"JobID": c.JOBID,
					"err":   err.Error(),
				}).Error("read error")
				return
			}
		}

		fdtdline := ParseStdOutLine(line)
		//exec hook function
		for _, fn := range c.onFDTDLineRecv {
			if fn != nil {
				fn(fdtdline)
			}
		}

		switch fdtdline.Type {
		case LineTypeComplete:
			c.setFDTDStatus(FDTDComplete)
			break loop
		case LineTypePlain:
			if fdtdline.PlainError == ErrorLicense {
				c.setFDTDStatus(FDTDFailed)
				return
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

func (c *FDTDSlurmTask) getJobInfo() (slurm.JobInfoMsg, error) {
	return slurm.LoadJobs(slurm.LoadJobsPayload{
		JobId: &c.JOBID,
	})
}
