# go-dict

Lazy go map value deambiguator.

# Installation

```shell
go get -u github.com/dcarbone/go-dict
```

# Usage

```go
package main

import (
	"fmt"

	"github.com/dcarbone/go-dict"
)

const (
	wellDocumentedJson = `{
    "definitelyString": false,
    "docsSayInt": "200",
    "billIsSureItsABool": "true"
}`
)

func main() {
	d, err := dict.UnmarshalJSON([]byte(wellDocumentedJson))
	if err != nil {
		panic(err.Error())
    }

	fmt.Printf("%[1]v (%[1]T)\n", d.MustGetString("definitelyString"))
	fmt.Printf("%[1]v (%[1]T)\n", d.MustGetInt("docsSayInt"))
	fmt.Printf("%[1]v (%[1]T)\n", d.MustGetBool("billIsSureItsABool"))
}
```
