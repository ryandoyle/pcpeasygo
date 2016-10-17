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
)

type pmValueAdapter interface {
	toUntypedMetric(value_format int, metric_type int, pm_value *pmapi.PmValue) (interface{}, error)
}

type pmValueAdapterImpl struct{
	pmapi pmapi.PMAPI
}

func (a pmValueAdapterImpl) toUntypedMetric(value_format int, metric_type int, pm_value *pmapi.PmValue) (interface{}, error) {
	pm_atom_value, err := a.pmapi.PmExtractValue(value_format, metric_type, pm_value)
	if(err != nil) {
		return nil, err
	}

	switch metric_type {
	case pmapi.PmType32:
		return pm_atom_value.Int32, nil
	case pmapi.PmTypeU32:
		return pm_atom_value.UInt32, nil
	case pmapi.PmType64:
		return pm_atom_value.Int64, nil
	case pmapi.PmTypeU64:
		return pm_atom_value.UInt64, nil
	case pmapi.PmTypeFloat:
		return pm_atom_value.Float, nil
	case pmapi.PmTypeDouble:
		return pm_atom_value.Double, nil
	case pmapi.PmTypeString:
		return pm_atom_value.String, nil
	}
	/* Shouldn't ever get here as PmExtractValue would exit earlier */
	return nil, nil
}