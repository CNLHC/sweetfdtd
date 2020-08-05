package slurm

/*
#cgo LDFLAGS: -L/usr/local/lib/slurm -lslurmfull
#include "root.h"
typedef char *chars;

#define SLUW_LIST(T, S)                               \
    T *sluw_alloc_##T(int s)                          \
    {                                                 \
        T *r = (T *)calloc(s + 1, sizeof(T));         \
        if (r)                                        \
            r[s] = S;                                 \
        return r;                                     \
    }                                                 \
    void sluw_set_##T(T *l, T v, int p) { l[p] = v; } \
    size_t sluw_len_##T(T *l)                         \
    {                                                 \
        size_t i = 0;                                 \
        if (l)                                        \
            while (l[i] != S)                         \
                i++;                                  \
        return i;                                     \
    }

SLUW_LIST(uint32_t, 0)
SLUW_LIST(int32_t, -1)
SLUW_LIST(chars, NULL)
*/
import (
	"C"
)

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	logging "github.com/sweetfdtd/pkg/log"
)

type Object struct {
	Type   string
	Offset unsafe.Pointer
}
type ObjectMap map[string]Object
type Table map[string]interface{}
type Request []byte

var logger = logging.GetLogger()

func sluw_get_name(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}

	str := strings.Split(s, "_")

	for i, v := range str {
		tmp := strings.SplitN(v, "", 2)
		tmp[0] = strings.ToUpper(tmp[0])
		str[i] = strings.Join(tmp, "")
	}

	return strings.Join(str, "")
}

func (t ObjectMap) Add(data interface{}) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Ptr {
		return
	}

	val = val.Elem()

	if !val.CanAddr() {
		return
	}

	ptr := val.Addr().Pointer()

	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		offset := val.Type().Field(i).Offset
		field := val.Field(i)
		if name == "_" {
			continue
		}
		t[sluw_get_name(name)] = Object{
			field.Type().String(),
			unsafe.Pointer(ptr + offset),
		}
	}
}

func (t ObjectMap) BindRequest(sreq Request, api_call func() (*Table, error)) (*Table, error) {
	var req map[string]*json.RawMessage
	err := json.Unmarshal(sreq, &req)
	if err != nil {
		return nil, err
	}

	logger.Debugf("BindRequest: Unmarshall result %+v", req)
	for key, value := range req {
		dst, ok := t[key]

		if !ok {
			errMsg := fmt.Sprintf("Invalid Field %s", key)
			logger.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		logger.Debugf("BindRequest: Reflected target meta %+v for %s", dst, key)
		var err error
		switch dst.Type {
		case "*slurm._Ctype_char":
			var s string
			err = json.Unmarshal(*value, &s)
			tmp := C.CString(s)
			logger.Debugf("BindRequest: Got value: %s", s)
			*(**C.char)(dst.Offset) = tmp
			defer C.free(unsafe.Pointer(tmp))

		case "slurm._Ctype_uint",
			"slurm._Ctype_uint32_t":
			var i uint32
			err = json.Unmarshal(*value, &i)
			*(*C.uint32_t)(dst.Offset) = C.uint32_t(i)
		case "slurm._Ctype_uint16_t":
			var i uint16
			err = json.Unmarshal(*value, &i)
			*(*C.uint16_t)(dst.Offset) = C.uint16_t(i)
		case "slurm._Ctype_uint8_t":
			var i uint8
			err = json.Unmarshal(*value, &i)
			*(*C.uint8_t)(dst.Offset) = C.uint8_t(i)
		case "slurm._Ctype_int32_t":
			var i int32
			err = json.Unmarshal(*value, &i)
			*(*C.int32_t)(dst.Offset) = C.int32_t(i)
		case "slurm._Ctype_int16_t":
			var i int16
			err = json.Unmarshal(*value, &i)
			*(*C.int16_t)(dst.Offset) = C.int16_t(i)
		case "slurm._Ctype_int8_t":
			var i int8
			err = json.Unmarshal(*value, &i)
			*(*C.int8_t)(dst.Offset) = C.int8_t(i)
		case "slurm._Ctype_time_t":
			var i uint64
			err = json.Unmarshal(*value, &i)
			*(*C.time_t)(dst.Offset) = C.time_t(i)
		case "*slurm._Ctype_uint32_t":
			var ai []uint32
			err = json.Unmarshal(*value, &ai)
			tmp := C.sluw_alloc_uint32_t(C.int(len(ai)))
			defer C.free(unsafe.Pointer(tmp))
			for i := 0; i < len(ai); i++ {
				C.sluw_set_uint32_t(tmp, C.uint32_t(ai[i]), C.int(i))
			}
			*(**C.uint32_t)(dst.Offset) = tmp
		case "*slurm._Ctype_int32_t":
			var ai []int32
			err = json.Unmarshal(*value, &ai)
			tmp := C.sluw_alloc_int32_t(C.int(len(ai)))
			defer C.free(unsafe.Pointer(tmp))
			for i := 0; i < len(ai); i++ {
				C.sluw_set_int32_t(tmp, C.int32_t(ai[i]), C.int(i))
			}
			*(**C.int32_t)(dst.Offset) = tmp
		case "**slurm._Ctype_char":
			var as []string
			err = json.Unmarshal(*value, &as)
			tmp := C.sluw_alloc_chars(C.int(len(as)))
			defer C.free(unsafe.Pointer(tmp))
			for i := 0; i < len(as); i++ {
				tmp2 := C.CString(as[i])
				defer C.free(unsafe.Pointer(tmp2))
				C.sluw_set_chars(tmp, tmp2, C.int(i))
			}
			*(***C.char)(dst.Offset) = (**C.char)(tmp)
		default:
			logger.Error(key, reflect.TypeOf(dst), "not supported")
		}
		if err != nil {
			errMsg := fmt.Sprintf("Bad value for key: ", key)
			return nil, errors.New(errMsg)
		}
	}
	if res, err := api_call(); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func GetRes(data interface{}) *Table {
	ret := make(Table)
	val := reflect.ValueOf(data).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		if f.Name == "_" {
			continue
		}
		name := sluw_get_name(f.Name)
		v := val.Field(i)
		switch f.Type.String() {
		case "slurm._Ctype_uint8_t",
			"slurm._Ctype_uint16_t",
			"slurm._Ctype_ushort",
			"slurm._Ctype_uint",
			"slurm._Ctype_ulong",
			"slurm._Ctype_uint32_t",
			"slurm._Ctype_uint64_t":
			ret[name] = uint(v.Uint())
		case "slurm._Ctype_int8_t",
			"slurm._Ctype_int16_t",
			"slurm._Ctype_int32_t",
			"slurm._Ctype_long",
			"slurm._Ctype_int64_t",
			"slurm._Ctype_time_t", // why not..
			"slurm._Ctype_int":
			ret[name] = int(v.Int())
		case "*slurm._Ctype_uint32_t":
			if v.Pointer() == 0 {
				ret[name] = make([]uint, 0)
				break
			}
			data := unsafe.Pointer(v.Pointer())
			count := int(C.sluw_len_uint32_t((*C.uint32_t)(data)))
			array := make([]uint, count)
			carray := *(*[]C.uint32_t)(unsafe.Pointer(&reflect.SliceHeader{
				Data: uintptr(data),
				Len:  count,
				Cap:  count,
			}))
			for k := 0; k < count; k++ {
				array[k] = uint(carray[k])
			}
			ret[name] = array
		case "*slurm._Ctype_int32_t":
			if v.Pointer() == 0 {
				ret[name] = make([]int, 0)
				break
			}
			data := unsafe.Pointer(v.Pointer())
			count := int(C.sluw_len_int32_t((*C.int32_t)(data)))
			array := make([]int, count)
			carray := *(*[]C.int32_t)(unsafe.Pointer(&reflect.SliceHeader{
				Data: uintptr(data),
				Len:  count,
				Cap:  count,
			}))
			for k := 0; k < count; k++ {
				array[k] = int(carray[k])
			}
			ret[name] = array
		case "*slurm._Ctype_char":
			if v.Pointer() == 0 {
				ret[name] = nil
				break
			}
			ret[name] = C.GoString((*C.char)(unsafe.Pointer(v.Pointer())))
		default:
		}
	}

	return &ret
}

func SlurmError() error {
	errno := C.slurm_get_errno()
	errno_str := "SLURM-" + strconv.Itoa(int(errno)) + " " + C.GoString(C.slurm_strerror(errno))
	return errors.New(errno_str)
}
