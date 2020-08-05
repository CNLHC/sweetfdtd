package fdtd

import (
	"fmt"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
)

type FDTDTaskSet struct {
	Tasks []FDTDSlurmTask
	wg    sync.WaitGroup
}

func (c *FDTDTaskSet) BuildFromPath(fp string) error {
	if files, err := ListAllFile(fp, ".fsp"); err != nil {
		return err
	} else {
		for _, file := range files {
			c.Tasks = append(c.Tasks,
				FDTDSlurmTask{
					FSPFile: path.Join(fp, file),
					Wg:      &c.wg,
				})
		}
	}
	return nil
}

func (c *FDTDTaskSet) Run() {
	logger.WithFields(logrus.Fields{
		"set": fmt.Sprintf("%+v", c.Tasks),
	}).Info("Run taskset")
	for _, task := range c.Tasks {
		c.wg.Add(1)
		task.Submit()
	}
	c.wg.Wait()
}

func (c *FDTDTaskSet) retry() {
	for _, task := range c.Tasks {
		if task.Status.FDTD == FDTDFailed {
			task.Submit()
		}
	}
}
