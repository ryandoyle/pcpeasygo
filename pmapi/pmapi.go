package pmapi
// #cgo LDFLAGS: -lpcp
// #include <pcp/pmapi.h>
import "C"
import (
	"unsafe"
	"errors"
	"runtime"
)

type PmapiContext struct {
	context int
}

func finalizer(c *PmapiContext) {
	C.pmDestroyContext(C.int(c.context))
}

func PmNewContext(hostname string) (*PmapiContext, error) {
	hostname_ptr := C.CString(hostname)
	defer C.free(unsafe.Pointer(hostname_ptr))

	context_id := int(C.pmNewContext(C.PM_CONTEXT_HOST, hostname_ptr))
	if (context_id < 0) {
		return nil, errors.New(PmErrStr(context_id))
	}

	context := &PmapiContext{
		context: context_id,
	}

	runtime.SetFinalizer(context, finalizer)

	return context, nil
}

func (c *PmapiContext) PmGetContextHostname() string {
	string_buffer := make([]C.char, C.MAXHOSTNAMELEN)
	raw_char_ptr := (*C.char)(unsafe.Pointer(&string_buffer[0]))

	C.pmGetContextHostName_r(C.int(c.context), raw_char_ptr, C.MAXHOSTNAMELEN)

	return C.GoString(raw_char_ptr)
}

func PmErrStr(error_no int) string {
	string_buffer := make([]C.char, C.PM_MAXERRMSGLEN)
	raw_char_ptr := (*C.char)(unsafe.Pointer(&string_buffer[0]))

	C.pmErrStr_r(C.int(error_no), raw_char_ptr, C.PM_MAXERRMSGLEN)

	return C.GoString(raw_char_ptr)
}

func (c *PmapiContext) GetContextId() int {
	return c.context
}