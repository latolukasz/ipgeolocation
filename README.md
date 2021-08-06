# IP geolocation 

## Importing data from csv files

You need to convert geolocation database *.csv files into *.db files that
are used by this library.

```go
package yourpackage

import "github.com/latolukasz/ipgeolocation"

func main() {
    err := ipgeolocation. Import("/path/to/csv/files/")
    if err != nil {
        panic(err)
    }
}
```

## Searching 

```go
package yourpackage

import (
	"github.com/latolukasz/ipgeolocation"
	"fmt"
)

func main() {
    err := ipgeolocation.InitDB("/path/to/db/files/") // run it only once in your app
    if err != nil {
        panic(err)
    }
    record, err := ipgeolocation.Search("37.143.210.32")
    if err != nil {
        panic(err)
    }
    if record != nil {
        fmt.Printf("country: %s, city: %s", record.Country.Name, record.City)
    } else {
        fmt.Printf("not found")
    }
}
```