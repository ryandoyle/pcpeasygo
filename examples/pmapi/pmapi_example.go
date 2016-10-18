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

package main

import (
	"fmt"
	"github.com/ryandoyle/pcpeasygo/pmapi"
	"runtime"
	"time"
)

func main() {
	doRun()
	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)
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

	vset := result.VSet[0]
	for _, pm_value := range vset.VList {
		val, _ := pmapi.PmExtractValue(vset.ValFmt, desc.Type, pm_value)
		instances_and_values[indoms[pm_value.Inst]] = val.Int32
	}

	fmt.Printf("context id is %v\nhostname is %v\nmetric is %v\npmdesc is %v\nindoms: %v\nvalues: %v\n",
		context.GetContextId(), host, metric, desc, indoms, instances_and_values)
}
