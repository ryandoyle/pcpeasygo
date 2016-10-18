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

	disk_all_read_metric, _   := a.Metric("disk.all.read")
	disk_partition_metrics, _ := a.Metrics("disk.partitions.read", "disk.partitions.write")

	fmt.Printf("metric(no instance): %+v\n", disk_all_read_metric)
	fmt.Printf("metrics(with instance): %+v\n", disk_partition_metrics)
}
```


For more, see `examples/` in the project root. API can change without notice as 
these bindings are being written.

## License
This project is licensed under the MIT license

PCP is licenced under LGPL 2.1. Its source code can be found at https://github.com/performancecopilot/pcp
