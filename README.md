# ifacegen

ifacegen is a code generation tool that generates interfaces from structs.

## Usage

A go generate comment can be added to the struct:


```go
package client

import (
	"context"
	"net/url"
)

//go:generate go run github.com/ChrisRx/ifacegen
type Client struct {
    ...
}


// ExportedMethod is an example method for Client that is exported.
func (c *Client) ExportedMethod(ctx context.Context, u *url.URL) (string, error) {
    ...
}

func (c *Client) unExportedMethod() error {
    ...
}
```

and a file will be produced in the same package named `zz_Client.iface.go` containing the generated interface:


```go
// Code generated by ifacegen. DO NOT EDIT.

package client

import (
	"context"
	"net/url"
)

type Interface interface {
	// ExportedMethod is an example method for Client that is exported.
	ExportedMethod(ctx context.Context, u *url.URL) (string, error)
}
```

ifacegen will determine the target struct based upon where in the code the go generate comment is placed.

Flags can be used specify the target struct (`--struct-name`) or the generated interface name (`--iface`):

```go
//go:generate go run github.com/ChrisRx/ifacegen --struct-name client -iface Client
```