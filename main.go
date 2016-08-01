package main

import (
	"fmt"
	"github.com/ryandoyle/pcpeasygo/pmapi"
)

func main() {
	context, err := pmapi.PmNewContext("localhost")
	if err != nil {
		panic(err)
	}
	fmt.Printf("context id is %v, hostname is %v\n", context.GetContextId(), context.PmGetContextHostname())
}
