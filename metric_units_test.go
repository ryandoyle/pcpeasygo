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

package pcpeasy

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"reflect"
)

var metricUnitsTests = []struct{
	desc string
	in  pmapi.PmDesc
	out metricInfo
}{
	{"counter semantics", pmapi.PmDesc{Sem:pmapi.PmSemCounter}, metricInfo{semantics:"counter", _type:reflect.Int32}},
	{"discrete semantics", pmapi.PmDesc{Sem:pmapi.PmSemDiscrete}, metricInfo{semantics:"discrete", _type:reflect.Int32}},
	{"instant semantics", pmapi.PmDesc{Sem:pmapi.PmSemInstant}, metricInfo{semantics:"instant", _type:reflect.Int32}},
	{"unknown semantics", pmapi.PmDesc{Sem:-123}, metricInfo{semantics:"unknown", _type:reflect.Int32}},

	{"bytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceByte}}, metricInfo{units:metricUnits{domain:"bytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"kilobytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceKByte}}, metricInfo{units:metricUnits{domain:"kilobytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"megabytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceMByte}}, metricInfo{units:metricUnits{domain:"megabytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"gigabytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceGByte}}, metricInfo{units:metricUnits{domain:"gigabytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"terabytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceTByte}}, metricInfo{units:metricUnits{domain:"terabytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"petabytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpacePByte}}, metricInfo{units:metricUnits{domain:"petabytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"exabytes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceEByte}}, metricInfo{units:metricUnits{domain:"exabytes"}, semantics:"unknown", _type:reflect.Int32}},
	{"count0 domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimCount:1,ScaleCount:0}}, metricInfo{units:metricUnits{domain:"count0"}, semantics:"unknown", _type:reflect.Int32}},
	{"count1 domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimCount:1,ScaleCount:1}}, metricInfo{units:metricUnits{domain:"count1"}, semantics:"unknown", _type:reflect.Int32}},
	{"count2 domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimCount:1,ScaleCount:2}}, metricInfo{units:metricUnits{domain:"count2"}, semantics:"unknown", _type:reflect.Int32}},
	{"nanoseconds domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeNSec}}, metricInfo{units:metricUnits{domain:"nanoseconds"}, semantics:"unknown", _type:reflect.Int32}},
	{"microseconds domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeUSec}}, metricInfo{units:metricUnits{domain:"microseconds"}, semantics:"unknown", _type:reflect.Int32}},
	{"milliseconds domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeMSec}}, metricInfo{units:metricUnits{domain:"milliseconds"}, semantics:"unknown", _type:reflect.Int32}},
	{"seconds domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeSec}}, metricInfo{units:metricUnits{domain:"seconds"}, semantics:"unknown", _type:reflect.Int32}},
	{"minutes domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeMin}}, metricInfo{units:metricUnits{domain:"minutes"}, semantics:"unknown", _type:reflect.Int32}},
	{"hours domain", pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeHour}}, metricInfo{units:metricUnits{domain:"hours"}, semantics:"unknown", _type:reflect.Int32}},

	/* Types */
	{"int32 type", pmapi.PmDesc{Type:pmapi.PmType32}, metricInfo{_type:reflect.Int32, semantics:"unknown"}},
	{"uint32 type", pmapi.PmDesc{Type:pmapi.PmTypeU32}, metricInfo{_type:reflect.Uint32, semantics:"unknown"}},
	{"int64 type", pmapi.PmDesc{Type:pmapi.PmType64}, metricInfo{_type:reflect.Int64, semantics:"unknown"}},
	{"uint64 type", pmapi.PmDesc{Type:pmapi.PmTypeU64}, metricInfo{_type:reflect.Uint64, semantics:"unknown"}},
	{"float type", pmapi.PmDesc{Type:pmapi.PmTypeFloat}, metricInfo{_type:reflect.Float32, semantics:"unknown"}},
	{"double type", pmapi.PmDesc{Type:pmapi.PmTypeDouble}, metricInfo{_type:reflect.Float64, semantics:"unknown"}},
	{"string type", pmapi.PmDesc{Type:pmapi.PmTypeString}, metricInfo{_type:reflect.String, semantics:"unknown"}},
	{"unknown type", pmapi.PmDesc{Type:pmapi.PmTypeEvent}, metricInfo{_type:reflect.Invalid, semantics:"unknown"}},

	/* Metrics with a range dimension */
	{
		"bytes/seconds",
	 	pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceByte, DimTime:-1, ScaleTime:pmapi.PmTimeSec}},
		metricInfo{units:metricUnits{domain:"bytes", _range:"seconds"}, semantics:"unknown", _type:reflect.Int32},
	},
	{
		"bytes/count3",
		pmapi.PmDesc{Units:pmapi.PmUnits{DimSpace:1,ScaleSpace:pmapi.PmSpaceByte, DimCount:-1, ScaleCount:3}},
		metricInfo{units:metricUnits{domain:"bytes", _range:"count3"}, semantics:"unknown", _type:reflect.Int32},
	},

	{
		"seconds/bytes",
		pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeSec, DimSpace:-1, ScaleSpace:pmapi.PmSpaceByte}},
		metricInfo{units:metricUnits{domain:"seconds", _range:"bytes"}, semantics:"unknown", _type:reflect.Int32},
	},
	{
		"seconds/count3",
		pmapi.PmDesc{Units:pmapi.PmUnits{DimTime:1,ScaleTime:pmapi.PmTimeSec, DimCount:-1, ScaleCount:3}},
		metricInfo{units:metricUnits{domain:"seconds", _range:"count3"}, semantics:"unknown", _type:reflect.Int32},
	},

	{
		"count3/bytes",
		pmapi.PmDesc{Units:pmapi.PmUnits{DimCount:1,ScaleCount:3, DimSpace:-1, ScaleSpace:pmapi.PmSpaceByte}},
		metricInfo{units:metricUnits{domain:"count3", _range:"bytes"}, semantics:"unknown", _type:reflect.Int32},
	},
	{
		"count3/seconds",
		pmapi.PmDesc{Units:pmapi.PmUnits{DimCount:1,ScaleCount:3, DimTime:-1, ScaleTime:pmapi.PmTimeSec}},
		metricInfo{units:metricUnits{domain:"count3", _range:"seconds"}, semantics:"unknown", _type:reflect.Int32},
	},

}

func TestToMetricUnits(t *testing.T) {
	adapter := pmDescAdapterImpl{}
	for test_number, test := range metricUnitsTests {
		assert.Equal(t, adapter.toMetricInfo(test.in), test.out, "test number: %v, description: \"%v\" ", test_number, test.desc)
	}
}