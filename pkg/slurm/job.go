package slurm

/*
#cgo LDFLAGS: -L/usr/local/lib/slurm -lslurmfull
#include "/root.h"
*/
import (
	"C"
)
import (
	"fmt"
	"reflect"
	"unsafe"
)

// func SubmitBatchJob(w http.ResponseWriter, r *http.Request) {
// 	var slreq C.job_desc_msg_t
// 	C.slurm_init_job_desc_msg(&slreq)

// 	obj := make(slurm.ObjectMap)
// 	obj.Add(&slreq)

// 	obj.Run(w, r, func() {
// 		var slres *C.submit_response_msg_t

// 		ret := C.slurm_submit_batch_job(&slreq, &slres)

// 		if ret != 0 {
// 			slurm.SlurmError(w, r)
// 			return
// 		}

// 		res := slurm.GetRes(slres)
// 		C.slurm_free_submit_response_response_msg(slres)

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(&res)
// 	})
// }

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

type LoadJobsPayload struct {
	UpdateTime C.time_t
	ShowFlags  C.uint16_t
}

func LoadJobs(Payload LoadJobsPayload) Table {
	var slres *C.job_info_msg_t

	ret := C.slurm_load_jobs(Payload.UpdateTime, &slres, Payload.ShowFlags)

	if ret != 0 {
		fmt.Printf("err")
		return nil
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

	C.slurm_free_job_info_msg(slres)
	return res

}
