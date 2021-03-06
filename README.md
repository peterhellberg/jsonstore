# jsonstore

[![Build Status](https://travis-ci.org/peterhellberg/jsonstore.svg?branch=master)](https://travis-ci.org/peterhellberg/jsonstore)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/jsonstore)](https://goreportcard.com/report/github.com/peterhellberg/jsonstore)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/jsonstore)

A client for the <https://www.jsonstore.io> API (backend source code is available at <https://github.com/bluzi/jsonstore>)

## Installation

    go get -u github.com/peterhellberg/jsonstore

## Usage example

```go
package main

import (
	"context"
	"fmt"

	"github.com/peterhellberg/jsonstore"
)

const secret = "3ba7860f742fc15d5b6e1508e2de1e0cde2c396f7c52a877905befb4e970eaaf"

type example struct {
	Number int
	Bool   bool
	String string
}

func main() {
	ctx := context.Background()

	store := jsonstore.New(jsonstore.Secret(secret))

	store.Post(ctx, "key", example{1234, true, "initial"})
	store.Put(ctx, "key/String", "modified")
	store.Delete(ctx, "key/Bool")

	var e example

	store.Get(ctx, "key", &e)

	fmt.Printf("%s -> %+v\n", store.URL("key"), e)
}
```

## License (MIT)

Copyright (c) 2018 [Peter Hellberg](https://c7.se/)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the "Software"),
> to deal in the Software without restriction, including without limitation
> the rights to use, copy, modify, merge, publish, distribute, sublicense,
> and/or sell copies of the Software, and to permit persons to whom the
> Software is furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included
> in all copies or substantial portions of the Software.

<img src="https://data.gopher.se/gopher/viking-gopher.svg" align="right" width="230" height="230">

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
> OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
> IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
> DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
> TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
> OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
