package slurm

/*
#cgo LDFLAGS: -L/usr/local/lib/slurm -lslurmfull
#include "root.h"
*/
import (
	"C"
)

type LoadJobsPayload struct {
	UpdateTime C.time_t
	ShowFlags  C.uint16_t
	JobId      *uint32
	Uid        *uint32
}

type JobInfoMsg struct {
	LastUpdate  uint64
	RecordCount uint64
	JobArray    []SlurmJobInfo
}

type SlurmJobInfo struct {
	StdErr         string
	StdIn          string
	StdOut         string
	UserName       string
	UserID         uint32
	WorkDir        string
	SubmitTime     uint64
	SuspendTime    uint64
	BatchHost      string
	Cluster        string
	Command        string
	Comment        string
	Dependency     string
	ExcNodes       string
	Features       string
	Gres           string
	Name           string
	Network        string
	Nodes          string
	Partition      string
	Qos            string
	ReqNodes       string
	ResvName       string
	SchedNodes     string
	StateDesc      string
	ThreadsPerCore uint16
}
