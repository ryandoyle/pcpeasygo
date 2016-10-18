> Work-in-progress Go bindings for [Performance Co-Pilot](http://pcp.io)

## Example
Using `pcpeasy` interface
```
package main

import (
	"fmt"
	"github.com/ryandoyle/pcpeasygo/pcpeasy"
)

func main() {
	a, _ := pcpeasy.NewAgent("localhost")

	metrics_with_no_instances, _ := a.Metrics("disk.all.read")
	metrics_with_instances, _ := a.Metrics("disk.partitions.read")

	fmt.Printf("metrics(no instance): %+v\n", metrics_with_no_instances)
	fmt.Printf("metrics(with instance): %+v\n", metrics_with_instances)
}
```


For more, see `examples/` in the project root. API can change without notice as 
these bindings are being written.

## License
This project is licensed under the MIT license

PCP is licenced under LGPL 2.1. Its source code can be found at https://github.com/performancecopilot/pcp
