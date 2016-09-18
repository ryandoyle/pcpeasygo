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
	metric, _ := context.PmLookupDesc(121634826)
	host, _ := context.PmGetContextHostname()
	fmt.Printf("context id is %v, hostname is %v, metric is %v\n", context.GetContextId(), host, metric)
}
