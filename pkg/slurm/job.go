package slurm

/*
#cgo LDFLAGS: -L/usr/local/lib/slurm -lslurmfull
#include "root.h"
*/
import (
	"C"
)
import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/mitchellh/mapstructure"
)

func SubmitBatchJob(sreq Request) (*Table, error) {
	var slreq C.job_desc_msg_t
	C.slurm_init_job_desc_msg(&slreq)

	obj := make(ObjectMap)
	obj.Add(&slreq)
	if res, err := obj.BindRequest(sreq, func() (*Table, error) {
		var slres *C.submit_response_msg_t
		ret := C.slurm_submit_batch_job(&slreq, &slres)
		if ret != 0 {
			return nil, SlurmError()
		}
		res := GetRes(slres)
		C.slurm_free_submit_response_response_msg(slres)
		return res, nil
	}); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// func NotifyJob(w http.ResponseWriter, r *http.Request) {
// 	opt := struct {
// 		job_id  C.uint32_t
// 		message *C.char
// 	}{}

// 	obj := make(slurm.ObjectMap)
// 	obj.Add(&opt)

// 	obj.Run(w, r, func() {
// 		ret := C.slurm_notify_job(opt.job_id, opt.message)

// 		if ret != 0 {
// 			slurm.SlurmError(w, r)
// 			return
// 		}
// 	})
// }

// func UpdateJob(w http.ResponseWriter, r *http.Request) {
// 	var slreq C.job_desc_msg_t
// 	C.slurm_init_job_desc_msg(&slreq)

// 	obj := make(slurm.ObjectMap)
// 	obj.Add(&slreq)

// 	obj.Run(w, r, func() {
// 		ret := C.slurm_update_job(&slreq)

// 		if ret != 0 {
// 			slurm.SlurmError(w, r)
// 			return
// 		}
// 	})
// }

func LoadJobs(Payload LoadJobsPayload) (JobInfoMsg, error) {
	var slres *C.job_info_msg_t
	var ret C.int

	if Payload.Uid != nil {
		ret = C.slurm_load_job_user(&slres, *(*C.uint)(Payload.Uid), Payload.ShowFlags)
	} else if Payload.JobId != nil {
		ret = C.slurm_load_job(&slres, *(*C.uint)(Payload.JobId), Payload.ShowFlags)
	} else {
		ret = C.slurm_load_jobs(Payload.UpdateTime, &slres, Payload.ShowFlags)
	}

	if ret != 0 {
		return JobInfoMsg{}, errors.New("ret is not 0")
	}

	data := unsafe.Pointer(slres.job_array)
	count := int(slres.record_count)
	carray := *(*[]C.job_info_t)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(data),
		Len:  count,
		Cap:  count,
	}))

	res := GetRes(slres)
	array := make([]*Table, count)
	for i := 0; i < count; i++ {
		array[i] = GetRes(&carray[i])
	}

	(*res)["JobArray"] = array

	var gores JobInfoMsg
	C.slurm_free_job_info_msg(slres)
	if err := mapstructure.Decode(res, &gores); err != nil {
		return JobInfoMsg{}, err
	} else {
		return gores, nil
	}

}
