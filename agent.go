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
	"reflect"
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"errors"
	"fmt"
)

type Metric struct {
	Name string
	Values []MetricValue
	Semantics string
	Type reflect.Kind
	Units Units
}

type Units struct {
	Domain string
	Range string
}

type MetricValue struct {
	Value interface{}
	Instance string
}

type agent struct {
	pmapi         pmapi.PMAPI
	pmDescAdapter pmDescAdapter
	pmValueAdapter pmValueAdapter
}

func NewAgent(host string) (*agent, error) {
	pmapi, err := pmapi.PmNewContext(pmapi.PmContextHost, host)
	if (err != nil) {
		return nil, err
	}
	return &agent{pmapi:pmapi, pmDescAdapter:pmDescAdapterImpl{}, pmValueAdapter:pmValueAdapterImpl{pmapi:pmapi}}, nil
}

func (a *agent) Metrics(metric_strings ...string) ([]Metric, error) {
	pmids, err := a.pmapi.PmLookupName(metric_strings...)
	if(err != nil) {
		return nil, err
	}
	/* Build up map of names and pmids */
	pmid_names := make(map[pmapi.PmID]string)
	for i, pmid := range pmids {
		pmid_names[pmid] = metric_strings[i]
	}

	pm_result, err := a.pmapi.PmFetch(pmids...)
	if(err != nil) {
		return nil, err
	}

	/* Just blow up if we can't get all the metrics we asked for */
	if(pm_result.NumPmID != len(metric_strings)) {
		return nil, errors.New("Error fetching all metrics")
	}

	metrics := make([]Metric, pm_result.NumPmID)
	for i, pm_value_set := range pm_result.VSet {
		metric, err := a.buildMetricFromPmValueSet(pm_value_set, pmid_names)
		if(err != nil) {
			return nil, err
		}
		metrics[i] = metric
	}


	return metrics, nil
}

func (a *agent) buildMetricFromPmValueSet(vset *pmapi.PmValueSet, pmid_names map[pmapi.PmID]string) (Metric, error) {
	metric_desc, err :=  a.pmapi.PmLookupDesc(vset.PmID)
	if(err != nil) {
		return Metric{}, err
	}
	metric_name := pmid_names[metric_desc.PmID]
	metric_info := a.pmDescAdapter.toMetricInfo(metric_desc)
	metric_values, err := a.buildMetricValues(vset, metric_desc)
	if(err != nil) {
		return Metric{}, err
	}
	return Metric{
		Type:metric_info._type,
		Name:metric_name,
		Semantics:metric_info.semantics,
		Units:Units{Domain:metric_info.units.domain, Range:metric_info.units._range},
		Values:metric_values,
	}, nil
}

func (a *agent) buildMetricValues(vset *pmapi.PmValueSet, metric_desc pmapi.PmDesc) ([]MetricValue, error) {
	if(vset.NumVal <= 0) {
		return nil, errors.New(fmt.Sprintf("metric \"%v\" contains no values or error \"%v\"", vset.PmID, vset.NumVal))
	}
	if(metric_desc.InDom == pmapi.PmInDomNull) {
		return a.buildMetricValuesForNullInstance(vset, metric_desc)
	} else {
		return a.buildMetricValuesForInstances(vset, metric_desc)
	}

}

func (a *agent) buildMetricValuesForNullInstance(vset *pmapi.PmValueSet, metric_desc pmapi.PmDesc) ([]MetricValue, error) {
	value, err := a.pmValueAdapter.toUntypedMetric(vset.ValFmt, metric_desc.Type, vset.VList[0])
	if (err != nil) {
		return nil, err
	}

	return []MetricValue{{
		Instance:"",
		Value:value,
	}}, nil
}

func (a *agent) buildMetricValuesForInstances(vset *pmapi.PmValueSet, metric_desc pmapi.PmDesc) ([]MetricValue, error) {
	ids_to_instance_names, err := a.pmapi.PmGetInDom(metric_desc.InDom)
	if(err != nil) {
		return nil, err
	}

	metric_values := make([]MetricValue, len(vset.VList))
	for i, pm_value := range vset.VList {
		value, err := a.pmValueAdapter.toUntypedMetric(vset.ValFmt, metric_desc.Type, pm_value)
		if(err != nil) {
			return nil, err
		}
		instance := ids_to_instance_names[pm_value.Inst]
		if(instance == "") {
			return nil, errors.New(fmt.Sprintf("Instance name for ID %v not found", pm_value.Inst))
		}
		metric_values[i] = MetricValue{Instance:instance, Value:value}
	}
	return metric_values, nil
}