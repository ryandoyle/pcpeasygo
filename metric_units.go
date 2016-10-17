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
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"fmt"
	"reflect"
)

type metricInfo struct {
	semantics string
	units metricUnits
	_type reflect.Kind
}

type metricUnits struct {
	domain string
	_range string
}

type pmDescAdapter interface {
	toMetricInfo(pm_desc pmapi.PmDesc) metricInfo
}

type pmDescAdapterImpl struct {}

func (a pmDescAdapterImpl) toMetricInfo(pm_desc pmapi.PmDesc) metricInfo {
	return metricInfo{
		semantics:semanticsString(pm_desc.Sem),
		units:createMetricUnits(pm_desc.Units),
		_type:createType(pm_desc.Type),
	}
}
func createType(_type int) reflect.Kind {
	switch _type {
	case pmapi.PmType32:
		return reflect.Int32
	case pmapi.PmTypeU32:
		return reflect.Uint32
	case pmapi.PmType64:
		return reflect.Int64
	case pmapi.PmTypeU64:
		return reflect.Uint64
	case pmapi.PmTypeFloat:
		return reflect.Float32
	case pmapi.PmTypeDouble:
		return reflect.Float64
	case pmapi.PmTypeString:
		return reflect.String
	}
	return reflect.Invalid
}

func createMetricUnits(units pmapi.PmUnits) metricUnits {
	if(units.DimSpace == 1 && units.DimTime == 0 && units.DimCount == 0) {
		return metricUnits{domain:spaceUnits(units)}
	} else if(units.DimSpace == 1 && units.DimTime == -1 && units.DimCount == 0) {
		return metricUnits{domain:spaceUnits(units), _range:timeUnits(units)}
	} else if(units.DimSpace == 1 && units.DimTime == 0 && units.DimCount == -1) {
		return metricUnits{domain:spaceUnits(units), _range:countUnits(units)}
	} else if(units.DimSpace == 0 && units.DimTime == 1 && units.DimCount == 0) {
		return metricUnits{domain:timeUnits(units)}
	} else if(units.DimSpace == -1 && units.DimTime == 1 && units.DimCount == 0) {
		return metricUnits{domain:timeUnits(units), _range:spaceUnits(units)}
	} else if(units.DimSpace == 0 && units.DimTime == 1 && units.DimCount == -1) {
		return metricUnits{domain:timeUnits(units), _range:countUnits(units)}
	} else if(units.DimSpace == 0 && units.DimTime == 0 && units.DimCount == 1) {
		return metricUnits{domain:countUnits(units)}
	} else if(units.DimSpace == -1 && units.DimTime == 0 && units.DimCount == 1) {
		return metricUnits{domain:countUnits(units), _range:spaceUnits(units)}
	} else if(units.DimSpace == 0 && units.DimTime == -1 && units.DimCount == 1) {
		return metricUnits{domain:countUnits(units), _range:timeUnits(units)}
	}
	return metricUnits{}
}

func countUnits(units pmapi.PmUnits) string {
	return fmt.Sprintf("count%v", units.ScaleCount)
}

func timeUnits(units pmapi.PmUnits) string {
	switch units.ScaleTime {
	case pmapi.PmTimeNSec:
		return "nanoseconds"
	case pmapi.PmTimeUSec:
		return "microseconds"
	case pmapi.PmTimeMSec:
		return "milliseconds"
	case pmapi.PmTimeSec:
		return "seconds"
	case pmapi.PmTimeMin:
		return "minutes"
	case pmapi.PmTimeHour:
		return "hours"
	}
	return "unknown time unit"
}


func spaceUnits(units pmapi.PmUnits) string {
	switch units.ScaleSpace {
	case pmapi.PmSpaceByte:
		return "bytes"
	case pmapi.PmSpaceKByte:
		return "kilobytes"
	case pmapi.PmSpaceMByte:
		return "megabytes"
	case pmapi.PmSpaceGByte:
		return "gigabytes"
	case pmapi.PmSpaceTByte:
		return "terabytes"
	case pmapi.PmSpacePByte:
		return "petabytes"
	case pmapi.PmSpaceEByte:
		return "exabytes"
	}
	return "unknown bytes"
}

func semanticsString(i int) string {
	switch i {
	case pmapi.PmSemCounter:
		return "counter"
	case pmapi.PmSemDiscrete:
		return "discrete"
	case pmapi.PmSemInstant:
		return "instant"
	}
	return "unknown"
}