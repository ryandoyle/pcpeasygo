package pmapi
// #cgo LDFLAGS: -lpcp
// #include <pcp/pmapi.h>
/*
// cgo does not support packed pmUnit struct. Define some helper functions
// to get the underlying data out of the struct
int getPmUnitsDimSpace(pmUnits units) {
	return units.dimSpace;
}
int getPmUnitsDimTime(pmUnits units) {
	return units.dimTime;
}
int getPmUnitsDimCount(pmUnits units) {
	return units.dimCount;
}

unsigned int getPmUnitsScaleSpace(pmUnits units) {
	return units.scaleSpace;
}
unsigned int getPmUnitsScaleTime(pmUnits units) {
	return units.scaleTime;
}
int getPmUnitsScaleCount(pmUnits units) {
	return units.scaleCount;
}

pmValueSet* getPmValueSetFromPmResult(int index, pmResult *pm_result) {
	return pm_result->vset[index];
}

pmValue getPmValueFromPmValueSet(int index, pmValueSet *pm_value_set) {
	return pm_value_set->vlist[index];
}

__int32_t getInt32FromPmAtomValue(pmAtomValue atom) {
	return atom.l;
}

__uint32_t getUInt32FromPmAtomValue(pmAtomValue atom) {
	return atom.ul;
}

__int64_t getInt64FromPmAtomValue(pmAtomValue atom) {
	return atom.ll;
}

__uint64_t getUInt64FromPmAtomValue(pmAtomValue atom) {
	return atom.ull;
}

float getFloatFromPmAtomValue(pmAtomValue atom) {
	return atom.f;
}

double getDoubleFromPmAtomValue(pmAtomValue atom) {
	return atom.d;
}

char *getStringFromPmAtomValue(pmAtomValue atom) {
	return atom.cp;
}

void freeStringFromPmAtomValue(pmAtomValue atom) {
	free(atom.cp);
}

void freePmValueBlockFromPmAtomValue(pmAtomValue atom) {
	free(atom.vbp);
}

*/
import "C"
import (
	"unsafe"
	"errors"
	"runtime"
	"time"
)

type PmapiContext struct {
	context int
}

type PmDesc struct {
	PmID PmID
	Type int
	InDom PmInDom
	Sem int
	Units PmUnits
}

type PmUnits struct {
	DimSpace int
	DimTime int
	DimCount int
	ScaleSpace uint
	ScaleTime uint
	ScaleCount int
}

type PmResult struct {
	cPmResult *C.pmResult
	vSet []*PmValueSet
}

type PmValueSet struct {
	cVset    *C.pmValueSet
	/* Hold reference to the underlying PmResult so it does not get GCed */
	pmResult *PmResult
	vList []*PmValue
}

type PmValue struct {
	cPmValue C.pmValue
	/* Hold reference to the underlying PmResult so it does not get GCed */
	pmResult *PmResult
}

type PmAtomValue struct {
	Int32 int32
	UInt32 uint32
	Int64 int64
	UInt64 uint64
	Float float32
	Double float64
	String string
}

type PmContextType int
type PmID uint
type PmInDom uint


const (
	PmContextHost = PmContextType(int(C.PM_CONTEXT_HOST))
	PmContextArchive = PmContextType(int(C.PM_CONTEXT_ARCHIVE))
	PmContextLocal = PmContextType(int(C.PM_CONTEXT_LOCAL))
	PmContextUndef = PmContextType(int(C.PM_CONTEXT_UNDEF))
	PmInDomNull = PmInDom(C.PM_INDOM_NULL)

	PmSpaceByte = uint(C.PM_SPACE_BYTE)
	PmSpaceKByte = uint(C.PM_SPACE_KBYTE)
	PmSpaceMByte = uint(C.PM_SPACE_MBYTE)
	PmSpaceGByte = uint(C.PM_SPACE_GBYTE)
	PmSpaceTByte = uint(C.PM_SPACE_TBYTE)
	PmSpacePByte = uint(C.PM_SPACE_PBYTE)
	PmSpaceEByte = uint(C.PM_SPACE_EBYTE)

	PmTimeNSec = uint(C.PM_TIME_NSEC)
	PmTimeUSec = uint(C.PM_TIME_USEC)
	PmTimeMSec = uint(C.PM_TIME_MSEC)
	PmTimeSec = uint(C.PM_TIME_SEC)
	PmTimeMin = uint(C.PM_TIME_MIN)
	PmTimeHour = uint(C.PM_TIME_HOUR)

	PmTypeNoSupport = int(C.PM_TYPE_NOSUPPORT)
	PmType32 = int(C.PM_TYPE_32)
	PmTypeU32 = int(C.PM_TYPE_U32)
	PmType64 = int(C.PM_TYPE_64)
	PmTypeU64 = int(C.PM_TYPE_U64)
	PmTypeFloat = int(C.PM_TYPE_FLOAT)
	PmTypeDouble = int(C.PM_TYPE_DOUBLE)
	PmTypeString = int(C.PM_TYPE_STRING)
	PmTypeAggregate= int(C.PM_TYPE_AGGREGATE)
	PmTypeAggregateStatic = int(C.PM_TYPE_AGGREGATE_STATIC)
	PmTypeEvent = int(C.PM_TYPE_EVENT)
	PmTypeHighResEvent = int(C.PM_TYPE_HIGHRES_EVENT)
	PmTypeUnknown = int(C.PM_TYPE_UNKNOWN)

	PmSemCounter = int(C.PM_SEM_COUNTER)
	PmSemInstant = int(C.PM_SEM_INSTANT)
	PmSemDiscrete = int(C.PM_SEM_DISCRETE)

	PmValInsitu = int(C.PM_VAL_INSITU)
	PmValDptr = int(C.PM_VAL_DPTR)
	PmValSptr = int(C.PM_VAL_SPTR)
)

func finalizer(c *PmapiContext) {
	C.pmDestroyContext(C.int(c.context))
}

func pmResultFinalizer(pm_result *PmResult) {
	C.pmFreeResult(pm_result.cPmResult)
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

func (c *PmapiContext) PmLookupDesc(pmid PmID) (PmDesc, error) {
	context_err := c.pmUseContext()
	if(context_err != nil) {
		return PmDesc{}, context_err
	}

	c_pmdesc := C.pmDesc{}

	err := int(C.pmLookupDesc(C.pmID(pmid), &c_pmdesc))
	if(err < 0) {
		return PmDesc{}, newPmError(err)
	}

	return PmDesc{
		PmID: PmID(c_pmdesc.pmid),
		Type: int(c_pmdesc._type),
		InDom: PmInDom(c_pmdesc.indom),
		Sem: int(c_pmdesc.sem),
		Units: PmUnits{
			DimSpace: int(C.getPmUnitsDimSpace(c_pmdesc.units)),
			DimTime: int(C.getPmUnitsDimTime(c_pmdesc.units)),
			DimCount: int(C.getPmUnitsDimCount(c_pmdesc.units)),
			ScaleSpace: uint(C.getPmUnitsScaleSpace(c_pmdesc.units)),
			ScaleTime: uint(C.getPmUnitsScaleTime(c_pmdesc.units)),
			ScaleCount: int(C.getPmUnitsScaleCount(c_pmdesc.units)),
		}}, nil

}

func (c *PmapiContext) PmGetInDom(indom PmInDom) (map[int]string, error) {
	context_err := c.pmUseContext()
	if(context_err != nil) {
		return nil, context_err
	}

	var c_instance_ids *C.int
	var c_instance_names **C.char

	err_or_number_of_instances := int(C.pmGetInDom(C.pmInDom(indom), &c_instance_ids, &c_instance_names))
	if(err_or_number_of_instances < 0) {
		return nil, newPmError(err_or_number_of_instances)
	}
	defer C.free(unsafe.Pointer(c_instance_ids))
	defer C.free(unsafe.Pointer(c_instance_names))

	/* Convert to a slice as we cannot do pointer arithmetic. As per
	   https://groups.google.com/forum/#!topic/golang-nuts/sV_f0VkjZTA */
	c_instance_ids_slice := (*[1 << 30]C.int)(unsafe.Pointer(c_instance_ids))
	c_instance_names_slice := (*[1 << 30]*C.char)(unsafe.Pointer(c_instance_names))

	indom_map := make(map[int]string)
	for i := 0; i < err_or_number_of_instances; i++ {
		indom_map[int(c_instance_ids_slice[i])] = C.GoString(c_instance_names_slice[i])
	}

	return indom_map, nil
}

func (c *PmapiContext) PmFetch(pmids ...PmID) (*PmResult, error) {
	number_of_pmids := len(pmids)

	var c_pm_result *C.pmResult
	c_pmids := (*C.pmID)(unsafe.Pointer(&pmids[0]))

	err := int(C.pmFetch(C.int(number_of_pmids), c_pmids, &c_pm_result))
	if(err < 0) {
		return &PmResult{}, newPmError(err)
	}

	result := &PmResult{cPmResult:c_pm_result}
	runtime.SetFinalizer(result, pmResultFinalizer)

	return result, nil
}

func PmExtractValue(value_format int, pm_type int, pm_value *PmValue) (PmAtomValue, error) {
	var c_pm_atom_value C.pmAtomValue

	err := int(C.pmExtractValue(C.int(value_format), &pm_value.cPmValue, C.int(pm_type), &c_pm_atom_value, C.int(pm_type)))
	if(err < 0) {
		return PmAtomValue{}, newPmError(err)
	}

	switch pm_type {
	case PmType32:
		return PmAtomValue{Int32:int32(C.getInt32FromPmAtomValue(c_pm_atom_value))}, nil
	case PmTypeU32:
		return PmAtomValue{UInt32:uint32(C.getUInt32FromPmAtomValue(c_pm_atom_value))}, nil
	case PmType64:
		return PmAtomValue{Int64:int64(C.getInt64FromPmAtomValue(c_pm_atom_value))}, nil
	case PmTypeU64:
		return PmAtomValue{UInt64:uint64(C.getUInt64FromPmAtomValue(c_pm_atom_value))}, nil
	case PmTypeFloat:
		return PmAtomValue{Float:float32(C.getFloatFromPmAtomValue(c_pm_atom_value))}, nil
	case PmTypeDouble:
		return PmAtomValue{Double:float64(C.getDoubleFromPmAtomValue(c_pm_atom_value))}, nil
	case PmTypeString:
		str := PmAtomValue{String:C.GoString(C.getStringFromPmAtomValue(c_pm_atom_value))}
		C.freeStringFromPmAtomValue(c_pm_atom_value)
		return str, nil
	case PmTypeAggregate:
	case PmTypeEvent:
	case PmTypeHighResEvent:
		C.freePmValueBlockFromPmAtomValue(c_pm_atom_value)
		return PmAtomValue{}, errors.New("Unsupported type")
	}
	return PmAtomValue{}, errors.New("Unknown type")
}

func (pm_result *PmResult) Timestamp() time.Time {
	return time.Unix(int64(pm_result.cPmResult.timestamp.tv_sec), int64(pm_result.cPmResult.timestamp.tv_usec) * 1000)
}

func (pm_result *PmResult) NumPmID() int {
	return int(pm_result.cPmResult.numpmid)
}

func (pm_result *PmResult) VSet() []*PmValueSet {
	/* Lazily cache the vSet */
	if(pm_result.vSet == nil) {
		vsets := make([]*PmValueSet, pm_result.NumPmID())
		for i := 0; i < pm_result.NumPmID(); i++ {
			vsets[i] = &PmValueSet{pmResult:pm_result, cVset:C.getPmValueSetFromPmResult(C.int(i), pm_result.cPmResult)}
		}
		pm_result.vSet = vsets
	}
	return pm_result.vSet
}

func (pm_value_set *PmValueSet) PmID() PmID {
	return PmID(pm_value_set.cVset.pmid)
}

func (pm_value_set *PmValueSet) NumVal() int {
	return int(pm_value_set.cVset.numval)
}

func (pm_value_set *PmValueSet) ValFmt() int {
	return int(pm_value_set.cVset.valfmt)
}

func (pm_value_set *PmValueSet) Vlist() []*PmValue {
	/* Lazily cache the vList */
	if(pm_value_set.vList == nil) {
		pm_values := make([]*PmValue, pm_value_set.NumVal())
		for i := 0; i < pm_value_set.NumVal(); i++ {
			pm_values[i] = &PmValue{
				pmResult:pm_value_set.pmResult,
				cPmValue:C.getPmValueFromPmValueSet(C.int(i), pm_value_set.cVset),
			}
		}
		pm_value_set.vList = pm_values
	}
	return pm_value_set.vList
}

func (pm_value *PmValue) Inst() int {
	return int(pm_value.cPmValue.inst)
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