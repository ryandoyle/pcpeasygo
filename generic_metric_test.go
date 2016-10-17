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
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"github.com/stretchr/testify/assert"
	"errors"
)

var pm_value = &pmapi.PmValue{}

func Test_ToUntypedMetric_forInt32(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{Int32:123}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmType32, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmType32, pm_value)

	assert.Equal(t, int32(123), value)
}

func Test_ToUntypedMetric_forUint32(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{UInt32:123}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeU32, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeU32, pm_value)

	assert.Equal(t, uint32(123), value)
}

func Test_ToUntypedMetric_forInt64(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{Int64:123}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmType64, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmType64, pm_value)

	assert.Equal(t, int64(123), value)
}

func Test_ToUntypedMetric_forUint64(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{UInt64:123}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeU64, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeU64, pm_value)

	assert.Equal(t, uint64(123), value)
}

func Test_ToUntypedMetric_forFloat(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{Float:123.456}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeFloat, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeFloat, pm_value)

	assert.Equal(t, float32(123.456), value)
}

func Test_ToUntypedMetric_forDouble(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{Double:123.456}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeDouble, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeDouble, pm_value)

	assert.Equal(t, float64(123.456), value)
}

func Test_ToUntypedMetric_forString(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}
	pm_atom_value := pmapi.PmAtomValue{String:"test"}

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeString, pm_value).Return(pm_atom_value, nil)

	value, _ := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeString, pm_value)

	assert.Equal(t, "test", value)
}

func Test_ToUntypedMetric_returnsAnErrorIfPmExtractValueReturnsAnError(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	adapter := pmValueAdapterImpl{pmapi:mock_pmapi}

	err := errors.New("error")

	mock_pmapi.On("PmExtractValue", pmapi.PmValDptr, pmapi.PmTypeString, pm_value).Return(nil, err)

	_, actual_err := adapter.toUntypedMetric(pmapi.PmValDptr, pmapi.PmTypeString, pm_value)

	assert.Equal(t, err, actual_err)
}