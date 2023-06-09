# go-athenahealth

go-athenahealth is an athenahealth API client for Go.

## Getting Started

### Install

```bash
$ go get github.com/eleanorhealth/go-athenahealth
```

### Basic Example

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eleanorhealth/go-athenahealth/athenahealth"
)

func main() {
	practiceID := "195900" // athenahealth shared sandbox practice ID.
	key := "your-api-key"
	secret := "your-api-secret"

	client := athenahealth.NewHTTPClient(&http.Client{}, practiceID, key, secret)

	p, err := client.GetPatient("1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s\n", p.FirstName, p.LastName)
}
```

### TokenCacher Example

Use `tokencacher.File` to cache API tokens to a file.

```go
client := athenahealth.NewHTTPClient(&http.Client{}, practiceID, key, secret).
    WithTokenCacher(tokencacher.NewFile("/tmp/athena_token.json"))
```
