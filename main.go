package main

import (
	"fmt"
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"runtime"
)

func main() {
	doRun()
	runtime.GC()
	runtime.GC()
	runtime.GC()
}

func doRun() {
	context, err := pmapi.PmNewContext(pmapi.PmContextHost, "localhost")
	if err != nil {
		panic(err)
	}
	host, _ := context.PmGetContextHostname()
	metric, _ := context.PmLookupName("sample.colour")
	desc, _ := context.PmLookupDesc(metric[0])
	indoms, _ := context.PmGetInDom(desc.InDom)
	result, _ := context.PmFetch(metric[0])

	instances_and_values := make(map[string]int32)

	vset := result.VSet()[0]
	for _, pm_value := range vset.Vlist() {
		val, _ := pmapi.PmExtractValue(vset.ValFmt(), desc.Type, pm_value)
		instances_and_values[indoms[pm_value.Inst()]] = val.Int32
	}

	fmt.Printf("context id is %v\nhostname is %v\nmetric is %v\npmdesc is %v\nindoms: %v\nvalues: %v\n",
		context.GetContextId(), host, metric, desc, indoms, instances_and_values)
}
