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
	"github.com/stretchr/testify/mock"
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"errors"
	"time"
	"reflect"
)

func TestNewAgent_createsANewAgentConnection(t *testing.T) {
	agent, _ := NewAgent("localhost")
	assert.NotNil(t, agent)
}

func TestNewAgent_returnsAnErrorForAnInvalidConnection(t *testing.T) {
	_, err := NewAgent("not-a-host")

	assert.Error(t, err)
}

type MockPMAPI struct {
	mock.Mock
}

type MockPmDescAdapter struct {
	mock.Mock
}

type MockPmValueAdapter struct {
	mock.Mock
}

func (m *MockPmValueAdapter) toUntypedMetric(value_format int, metric_type int, pm_value *pmapi.PmValue) (interface{}, error) {
	args := m.Called(value_format, metric_type, pm_value)
	value := args.Get(0)
	err := args.Error(1)
	if(value == nil) {
		return nil, err
	}
	return value.(interface{}), err
}

func (m *MockPMAPI) PmLookupName(names ...string) ([]pmapi.PmID, error) {
	args := m.Called(names)
	pmids := args.Get(0)
	err := args.Error(1)
	if(pmids == nil) {
		return nil, err
	}
	return pmids.([]pmapi.PmID), err
}

func (m *MockPMAPI) PmFetch(pmids ...pmapi.PmID) (*pmapi.PmResult, error) {
	args := m.Called(pmids)
	pm_result := args.Get(0)
	err := args.Error(1)
	if(pm_result == nil) {
		return nil, err
	}
	return pm_result.(*pmapi.PmResult), err
}

func (m *MockPMAPI) PmLookupDesc(pmid pmapi.PmID) (pmapi.PmDesc, error) {
	args := m.Called(pmid)
	pm_desc := args.Get(0)
	err := args.Error(1)
	if(pm_desc == nil) {
		return pmapi.PmDesc{}, err
	}
	return pm_desc.(pmapi.PmDesc), err
}

func (m *MockPMAPI) PmExtractValue(value_format int, pm_type int, pm_value *pmapi.PmValue) (pmapi.PmAtomValue, error) {
	args := m.Called(value_format, pm_type, pm_value)
	pm_atom_value := args.Get(0)
	err := args.Error(1)
	if(pm_atom_value == nil) {
		return pmapi.PmAtomValue{}, err
	}
	return pm_atom_value.(pmapi.PmAtomValue), err
}

func (m *MockPMAPI) PmGetInDom(indom pmapi.PmInDom) (map[int]string, error) {
	args := m.Called(indom)
	instances_to_names := args.Get(0)
	err := args.Error(1)
	if(instances_to_names == nil) {
		return nil, err
	}
	return instances_to_names.(map[int]string), err

}

func (m *MockPmDescAdapter) toMetricInfo(pm_desc pmapi.PmDesc) metricInfo {
	args := m.Called(pm_desc)
	return args.Get(0).(metricInfo)
}

func TestAgent_Metric_returnsAnErrorIfTheNameOfTheMetricCannotBeLookedUp(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	agent := &agent{pmapi:mock_pmapi}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return(nil, errors.New("PmLookupName error"))

	_, err := agent.Metrics("my.metric")

	assert.EqualError(t, err, "PmLookupName error")
}

func TestAgent_Metrics_returnsAnErrorIfAPmFetchFails(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	agent := &agent{pmapi:mock_pmapi}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return([]pmapi.PmID{123}, nil)
	mock_pmapi.On("PmFetch", []pmapi.PmID{123}).Return(nil, errors.New("PmFetch error"))

	_, err := agent.Metrics("my.metric")

	assert.EqualError(t, err, "PmFetch error")
}

func TestAgent_Metrics_returnsAnErrorIfTheNumberOfPMIDsIsLessThanAskedFor(t *testing.T) {
	pmIds := []pmapi.PmID{123}
	mock_pmapi := &MockPMAPI{}
	pm_value := &pmapi.PmResult{NumPmID:0}
	agent := &agent{pmapi:mock_pmapi}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return(pmIds, nil)
	mock_pmapi.On("PmFetch", pmIds).Return(pm_value, nil)

	_, err := agent.Metrics("my.metric")

	assert.EqualError(t, err, "Error fetching all metrics")
}

func TestAgent_Metrics_returnsAnErrorIfFetchingPmDescFails(t *testing.T) {
	pmIds := []pmapi.PmID{123}
	mock_pmapi := &MockPMAPI{}
	pm_value := &pmapi.PmResult{NumPmID:1, VSet:[]*pmapi.PmValueSet{{PmID:pmapi.PmID(123)}}}
	agent := &agent{pmapi:mock_pmapi}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return(pmIds, nil)
	mock_pmapi.On("PmFetch", pmIds).Return(pm_value, nil)
	mock_pmapi.On("PmLookupDesc", pmapi.PmID(123)).Return(nil, errors.New("pmdesc error"))

	_, err := agent.Metrics("my.metric")

	assert.EqualError(t, err, "pmdesc error")

}

func TestAgent_Metrics_returnsAMetricResult(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	mock_pmdesc_adapter := &MockPmDescAdapter{}
	mock_pmvalue_adapter := &MockPmValueAdapter{}
	agent := &agent{pmapi:mock_pmapi, pmDescAdapter:mock_pmdesc_adapter, pmValueAdapter:mock_pmvalue_adapter}
	pm_value := &pmapi.PmValue{Inst:pmapi.PmInNull}
	pm_result := &pmapi.PmResult{
		NumPmID:1,
		Timestamp:time.Unix(123,456),
		VSet:[]*pmapi.PmValueSet{{
			NumVal:1,
			PmID:pmapi.PmID(123),
			ValFmt:pmapi.PmValDptr,
			VList:[]*pmapi.PmValue{pm_value},
		}},
	}
	pm_desc := pmapi.PmDesc{Type:pmapi.PmType64, InDom:pmapi.PmInDomNull, PmID:pmapi.PmID(123)}
	metric_info := metricInfo{_type:reflect.Int64, semantics:"counter", units:metricUnits{_range:"seconds", domain:"megabytes"}}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return([]pmapi.PmID{123}, nil)
	mock_pmapi.On("PmFetch", []pmapi.PmID{123}).Return(pm_result, nil)
	mock_pmapi.On("PmLookupDesc", pmapi.PmID(123)).Return(pm_desc, nil)
	mock_pmdesc_adapter.On("toMetricInfo", pm_desc).Return(metric_info)
	mock_pmvalue_adapter.On("toUntypedMetric", pmapi.PmValDptr, pmapi.PmType64, pm_value).Return(int64(222), nil)

	actual_metrics, err := agent.Metrics("my.metric")

	expected_metrics := []Metric{{
		Name: "my.metric",
		Semantics: "counter",
		Type: reflect.Int64,
		Units: Units{Domain:"megabytes", Range: "seconds"},
		Values: []MetricValue{{Instance: "", Value: int64(222)}},
	}}

	assert.NoError(t, err)
	assert.Equal(t, expected_metrics, actual_metrics)
}

func TestAgent_Metrics_returnsAMetricResultForAMetricWithInstances(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	mock_pmdesc_adapter := &MockPmDescAdapter{}
	mock_pmvalue_adapter := &MockPmValueAdapter{}
	agent := &agent{pmapi:mock_pmapi, pmDescAdapter:mock_pmdesc_adapter, pmValueAdapter:mock_pmvalue_adapter}

	indom := pmapi.PmInDom(555)
	instance_1 := 111
	instance_2 := 222
	instance_names := map[int]string{instance_1:"inst1", instance_2:"inst2"}
	pm_value_1 := &pmapi.PmValue{Inst:instance_1}
	pm_value_2 := &pmapi.PmValue{Inst:instance_2}
	pmid := pmapi.PmID(123)
	pm_result := &pmapi.PmResult{
		NumPmID:1,
		Timestamp:time.Unix(123,456),
		VSet:[]*pmapi.PmValueSet{{
			NumVal:2,
			PmID:pmid,
			ValFmt:pmapi.PmValDptr,
			VList:[]*pmapi.PmValue{pm_value_1, pm_value_2},
		}},
	}
	pm_desc := pmapi.PmDesc{Type:pmapi.PmType64, InDom:indom, PmID:pmid}
	metric_info := metricInfo{_type:reflect.Int64, semantics:"counter", units:metricUnits{_range:"seconds", domain:"megabytes"}}

	mock_pmapi.On("PmLookupName", []string{"my.metric"}).Return([]pmapi.PmID{123}, nil)
	mock_pmapi.On("PmFetch", []pmapi.PmID{123}).Return(pm_result, nil)
	mock_pmapi.On("PmLookupDesc", pmid).Return(pm_desc, nil)
	mock_pmapi.On("PmGetInDom", indom).Return(instance_names, nil)
	mock_pmdesc_adapter.On("toMetricInfo", pm_desc).Return(metric_info)
	mock_pmvalue_adapter.On("toUntypedMetric", pmapi.PmValDptr, pmapi.PmType64, pm_value_1).Return(int64(881), nil)
	mock_pmvalue_adapter.On("toUntypedMetric", pmapi.PmValDptr, pmapi.PmType64, pm_value_2).Return(int64(882), nil)

	actual_metrics, err := agent.Metrics("my.metric")

	expected_metrics := []Metric{{
		Name: "my.metric",
		Semantics: "counter",
		Type: reflect.Int64,
		Units: Units{Domain:"megabytes", Range: "seconds"},
		Values: []MetricValue{
			{Instance: "inst1", Value: int64(881)},
			{Instance: "inst2", Value: int64(882)},
		},
	}}

	assert.NoError(t, err)
	assert.Equal(t, expected_metrics, actual_metrics)
}

func TestAgent_Metrics_returnsAErrorIfTheMetricContainsAnErrorEncodedInNumVal(t *testing.T) {
	mock_pmapi := &MockPMAPI{}
	mock_pmdesc_adapter := &MockPmDescAdapter{}
	mock_pmvalue_adapter := &MockPmValueAdapter{}
	agent := &agent{pmapi:mock_pmapi, pmDescAdapter:mock_pmdesc_adapter, pmValueAdapter:mock_pmvalue_adapter}

	metric_name := "my.metric"
	pmid := pmapi.PmID(123)
	pmids := []pmapi.PmID{pmid}
	error_encoded_in_numval := -12345
	pm_result := &pmapi.PmResult{
		NumPmID:1,
		VSet:[]*pmapi.PmValueSet{{
			NumVal:error_encoded_in_numval,
			PmID:pmid,
			ValFmt:pmapi.PmValDptr,
			VList:[]*pmapi.PmValue{},
		}},
	}
	pm_desc := pmapi.PmDesc{}
	metric_info := metricInfo{}

	mock_pmapi.On("PmLookupName", []string{metric_name}).Return(pmids, nil)
	mock_pmapi.On("PmFetch", pmids).Return(pm_result, nil)
	mock_pmapi.On("PmLookupDesc", pmid).Return(pm_desc, nil)
	mock_pmdesc_adapter.On("toMetricInfo", pm_desc).Return(metric_info)

	actual_metrics, err := agent.Metrics(metric_name)

	assert.Nil(t, actual_metrics)
	assert.EqualError(t, err, "metric \"123\" contains no values or error \"-12345\"")
}