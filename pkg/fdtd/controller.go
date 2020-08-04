package fdtd

import "path"

type FDTDTaskSet struct {
	Tasks []FDTDSlurmTask
}

func (c *FDTDTaskSet) BuildFromPath(fp string) error {
	if files, err := ListAllFile(fp, ".fsp"); err != nil {
		return err
	} else {
		for _, file := range files {
			c.Tasks = append(c.Tasks,
				FDTDSlurmTask{
					FSPFile: path.Join(fp, file),
				})
		}
	}
	return nil
}

func (c *FDTDTaskSet) Run() {
	for _, task := range c.Tasks {
		task.Submit()
	}
}

func (c *FDTDTaskSet) retry() {
	for _, task := range c.Tasks {
		if task.Status.FDTD == FDTDFailed {
			task.Submit()
		}
	}
}
