package slurm

import "testing"

func TestCommandBuildTest(t *testing.T) {
	tsk := &FDTDSlurmTask{}
	t.Logf(tsk.Submit())

}
