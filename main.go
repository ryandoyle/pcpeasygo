package main

import (
	"fmt"
	"github.com/ryandoyle/pcpeasygo/pmapi"
)

func main() {
	context, err := pmapi.PmNewContext(pmapi.PmContextHost, "localhost")
	if err != nil {
		panic(err)
	}
	host, _ := context.PmGetContextHostname()
	metric, _ := context.PmLookupName("sample.colour")
	desc, _ := context.PmLookupDesc(metric[0])
	indoms, _ := context.PmGetInDom(desc.InDom)

	fmt.Printf("context id is %v, hostname is %v, metric is %v, pmdesc is %v, indoms: %v\n", context.GetContextId(), host, metric,
		desc, indoms)
}
