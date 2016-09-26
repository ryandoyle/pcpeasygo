//Copyright (c) 2016 Ryan Doyle
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package pmapi

import (
	"testing"
	"reflect"
	"time"
)

var sampleDoubleMillionPmID PmID = 121634844
var sampleMillisecondsPmID PmID = 121634819
var sampleColourInDom PmInDom = 121634817
var sampleStringHulloPmID PmID = 121634847

func TestPmapiContext_PmGetContextHostname(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	hostname, _ := c.PmGetContextHostname()

	assertEquals(t, hostname, "ryandesktop")
}

func TestPmapiContext_PmLookupNameForASingleName(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	pmids, _ := c.PmLookupName("sample.double.million")

	assertEquals(t, pmids[0], sampleDoubleMillionPmID)
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
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillionPmID)

	assertEquals(t, pmdesc.PmID, sampleDoubleMillionPmID)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectType(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillionPmID)

	assertEquals(t, pmdesc.Type, PmTypeDouble)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectInDom(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillionPmID)

	assertEquals(t, pmdesc.InDom, PmInDomNull)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectSemantics(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleDoubleMillionPmID)

	assertEquals(t, pmdesc.Sem, PmSemInstant)
}

func TestPmapiContext_PmLookupDesc_HasTheCorrectUnits(t *testing.T) {
	pmdesc, _ := localContext().PmLookupDesc(sampleMillisecondsPmID)
	expected_units := PmUnits{
		DimTime: 1,
		ScaleTime: PmTimeMSec,
	}

	assertEquals(t, pmdesc.Units, expected_units)
}

func TestPmapiContext_PmLookupInDom_ReturnsTheInstanceNameMapping(t *testing.T) {
	indom, _ := localContext().PmGetInDom(sampleColourInDom)

	assertEquals(t, indom[0], "red")
	assertEquals(t, indom[1], "green")
	assertEquals(t, indom[2], "blue")
}

func TestPmapiContext_PmGetInDom_ReturnsANilErrorForValidInDoms(t *testing.T) {
	_, err := localContext().PmGetInDom(sampleColourInDom)

	assertNil(t, err)
}

func TestPmapiContext_PmGetInDom_ReturnsAnErrorForIncorrectInDoms(t *testing.T) {
	_, err := localContext().PmGetInDom(PmInDom(123))

	assertNotNil(t, err)
}

func TestPmapiContext_PmFetch_returnsAPmResultWithATimestamp(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertWithinDuration(t, pm_result.Timestamp, time.Now(), time.Second)
}

func TestPmapiContext_PmFetch_returnsAPmResultWithTheNumberOfPMIDs(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertEquals(t, pm_result.NumPmID, 1)
}

func TestPmapiContext_PmFetch_returnsAVSet_withAPmID(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertEquals(t, pm_result.VSet[0].PmID, sampleDoubleMillionPmID)
}

func TestPmapiContext_PmFetch_returnsAVSet_withNumval(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertEquals(t, pm_result.VSet[0].NumVal, 1)
}

func TestPmapiContext_PmFetch_returnsAVSet_withValFmt(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertEquals(t, pm_result.VSet[0].ValFmt, PmValDptr)
}

func TestPmapiContext_PmFetch_returnsAVSet_withVlist_withAPmValue_withAnInst(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	assertEquals(t, pm_result.VSet[0].VList[0].Inst, -1)
}

func TestPmExtractValue_forADoubleValue(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleDoubleMillionPmID)

	value := pm_result.VSet[0].VList[0]

	atom, _ := PmExtractValue(PmValDptr, PmTypeDouble, value)

	assertEquals(t, atom.Double, 1000000.0)
}

func TestPmExtractValue_forAStringValue(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleStringHulloPmID)

	value := pm_result.VSet[0].VList[0]

	atom, _ := PmExtractValue(PmValDptr, PmTypeString, value)

	assertEquals(t, atom.String, "hullo world!")
}

func TestPmExtractValue_returnsAnErrorWhenTryingToExtractTheWrongType(t *testing.T) {
	pm_result, _ := localContext().PmFetch(sampleStringHulloPmID)

	value := pm_result.VSet[0].VList[0]

	_, err := PmExtractValue(PmValDptr, PmType64, value)

	assertNotNil(t, err)
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

func assertWithinDuration(t *testing.T, time1 time.Time, time2 time.Time, duration time.Duration) {
	rounded1 := time1.Round(duration)
	rounded2 := time2.Round(duration)
	if(!rounded1.Equal(rounded2)) {
		t.Errorf("Expected time: %v and time: %v to be within %v of each other", time1, time2, duration)
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