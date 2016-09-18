package pmapi

import (
	"testing"
	"reflect"
)

var sampleDoubleMillion PmID = 121634844
var sampleMilliseconds PmID = 121634819

func TestPmapiContext_PmGetContextHostname(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	hostname, _ := c.PmGetContextHostname()

	assertEquals(t, hostname, "ryandesktop")
}

func TestPmapiContext_PmLookupNameForASingleName(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	pmids, _ := c.PmLookupName("sample.double.million")

	assertEquals(t, pmids[0], sampleDoubleMillion)
}

func TestPmapiContext_PmLookupNameReturnsAnErrorForUnknownNames(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	_, err := c.PmLookupName("not.a.name")

	assertNotNil(t, err)
}

func TestPmapiContext_PmLookupNameForMultipleNames(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	pmids, _ := c.PmLookupName("sample.long.one", "sample.ulong.hundred",)

	assertEquals(t, pmids[1], PmID(121634911))
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectPmID(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillion)

	assertEquals(t, pmdesc.PmID, sampleDoubleMillion)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectType(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillion)

	assertEquals(t, pmdesc.Type, PmTypeDouble)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectInDom(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillion)

	assertEquals(t, pmdesc.InDom, PmInDomNull)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectSemantics(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillion)

	assertEquals(t, pmdesc.Sem, PmSemInstant)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectUnits(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleMilliseconds)
	expected_units := PmUnits{
		DimTime: 1,
		ScaleTime: PmTimeMSec,
	}

	assertEquals(t, pmdesc.Units, expected_units)
}

func TestPmNewContext_withAnInvalidHostHasANilContext(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "not-a-host")

	assertNil(t, c)
}

func TestPmNewContext_withAnInvalidHostHasAnError(t *testing.T) {
	_, err := PmNewContext(PmContextHost, "not-a-host")

	assertNotNil(t, err)
}

func TestPmNewContext_hasAValidContext(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")

	assertNotNil(t, c)
}

func TestPmNewContext_hasANilError(t *testing.T) {
	_, err := PmNewContext(PmContextHost, "localhost")

	assertNil(t, err)
}

func TestPmNewContext_supportsALocalContext(t *testing.T) {
	c, _ := PmNewContext(PmContextLocal, "")

	assertNotNil(t, c)
}

func assertEquals(t *testing.T, a interface{}, b interface{}) {
	if(a != b) {
		t.Errorf("expected %v, got %v", b, a)
	}
}

func assertNil(t *testing.T, a interface{}) {
	if ! isNil(a) {
		t.Errorf("expected nil but got %v", a)
	}
}

func assertNotNil(t *testing.T, a interface{}) {
	if isNil(a) {
		t.Errorf("expected not nil but got %v", a)
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	return reflect.ValueOf(object).IsNil()
}

func localContext() *PmapiContext {
	c, _ := PmNewContext(PmContextHost, "localhost")
	return c
}