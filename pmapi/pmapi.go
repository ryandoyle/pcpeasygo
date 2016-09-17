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

type PmContextType int
type PmID uint

const (
	PmContextHost = PmContextType(int(C.PM_CONTEXT_HOST))
	PmContextArchive = PmContextType(int(C.PM_CONTEXT_ARCHIVE))
	PmContextLocal = PmContextType(int(C.PM_CONTEXT_LOCAL))
	PmContextUndef = PmContextType(int(C.PM_CONTEXT_UNDEF))
)

func finalizer(c *PmapiContext) {
	C.pmDestroyContext(C.int(c.context))
}

func PmNewContext(context_type PmContextType, host_or_archive string) (*PmapiContext, error) {
	host_or_archive_ptr := C.CString(host_or_archive)
	defer C.free(unsafe.Pointer(host_or_archive_ptr))

	context_id := int(C.pmNewContext(C.int(context_type), host_or_archive_ptr))
	if (context_id < 0) {
		return nil, newPmError(context_id)
	}

	context := &PmapiContext{
		context: context_id,
	}

	runtime.SetFinalizer(context, finalizer)

	return context, nil
}

func (c *PmapiContext) PmGetContextHostname() (string, error) {
	err := c.pmUseContext()
	if(err != nil) {
		return "", err
	}
	string_buffer := make([]C.char, C.MAXHOSTNAMELEN)
	raw_char_ptr := (*C.char)(unsafe.Pointer(&string_buffer[0]))

	C.pmGetContextHostName_r(C.int(c.context), raw_char_ptr, C.MAXHOSTNAMELEN)

	return C.GoString(raw_char_ptr), nil
}

func (c *PmapiContext) PmLookupName(names ...string) ([]PmID, error) {
	context_err := c.pmUseContext()
	if(context_err != nil) {
		return nil, context_err
	}

	number_of_names := len(names)
	c_pmids := make([]C.pmID, number_of_names)
	c_names := make([]*C.char, number_of_names)

	/* Build c_names as copies of the original names */
	for i, name := range names {
		name_ptr := C.CString(name)
		c_names[i] = name_ptr
		defer C.free(unsafe.Pointer(name_ptr))
	}

	/* Do the actual lookup */
	err := int(C.pmLookupName(C.int(number_of_names), &c_names[0], &c_pmids[0]))
	if(err < 0 ) {
		return nil, newPmError(err)
	}

	/* Collect up the C.pmIDs into Go PmID's. Originally when returning the slice that was passed
	into pmLookupName was resulting in bit length errors between Go's uint and C unsigned int */
	pmids := make([]PmID, number_of_names)
	for i, c_pmid := range c_pmids {
		pmids[i] = PmID(c_pmid)
	}

	return pmids, nil
}


func (c *PmapiContext) pmUseContext() error {
	err := int(C.pmUseContext(C.int(c.context)))
	if(err < 0) {
		return newPmError(err)
	}
	return nil
}

func newPmError(err int) error {
	return errors.New(pmErrStr(err))
}

func pmErrStr(error_no int) string {
	string_buffer := make([]C.char, C.PM_MAXERRMSGLEN)
	raw_char_ptr := (*C.char)(unsafe.Pointer(&string_buffer[0]))

	C.pmErrStr_r(C.int(error_no), raw_char_ptr, C.PM_MAXERRMSGLEN)

	return C.GoString(raw_char_ptr)
}

func (c *PmapiContext) GetContextId() int {
	return c.context
}