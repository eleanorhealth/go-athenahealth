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

## X-Request-Id

Clients can obtain the X-Request-Id sent on the request to athena from the
`X-Request-Id` header on the `http.Request` for use in tracing, logging, and debugging.

```go
xRequestId := req.Header.Get(athenahealth.XRequestIDHeaderKey))
```

See the athena [Best Practices](https://docs.athenahealth.com/api/guides/best-practices) guide for more details about X-Request-Id and other recommended  practices.

## Method Signatures Required vs. Optional Fields

All methods that perform network or filesystem IO will accept a context for idiomatic propagation.

While not yet consistently applied across this repo, required fields should be surfaced as top-level arguments in method signatures (including path parameters, query parameters, and request body fields) while optional fields should be embedded as pointer fields in an optional pointer struct. The aim is to provide a consistent and natural pattern for interacting with this client interface as well as to add compile-time safety to required fields.

### Required and Optional Fields

Optional fields options struct should by the last argument in the method signature, preceded by top-level required fields.

```go
type CreatePatientOptions struct {
	Address1              *string
	Address2              *string
	City                  *string
	Email                 *string
	HomePhone             *string
	MiddleName            *string
	MobilePhone           *string
	Notes                 *string
	Sex                   *string
	SSN                   *string
	State                 *string
	Status                *string
	Zip                   *string
	BypassPatientMatching *bool
}

type Client interface {
	// ...
	CreatePatient(ctx context.Context, departmentID, firstName, lastName string, dob time.Time, opts *CreatePatientOptions) (string, error)
	// ...
}
```

> [samber/lo#toptr](https://github.com/samber/lo#toptr) or a custom rolled to-pointer helper is recommended to improve the ergonomics of passing the pointers
> ```go
> func ToPtr[T any](x T) *T { return &x }
> ``` 

### Only Required Fields

If there are only required fields, options struct is omitted.

```go
type Client interface {
	// ...
	SearchAllergies(ctx context.Context, searchVal string) ([]*Allergy, error)
	// ...
}
```

### Only Optional Fields

If there are only optional fields, expect a standalone options struct.

```go
type Client interface {
	// ...
	ListChangedPatients(ctx context.Context, opts *ListChangedPatientOptions) ([]*Patient, error)
	// ...
}
```

### Neither Required or Optional Fields

If neither required or optional fields are present, expect neither.

```go
type Client interface {
	// ...
	ListAppointmentCustomFields(context.Context) ([]*AppointmentCustomField, error)
	// ...
}
```
