/*

Package jsonstore is a client for the www.jsonstore.io API

Installation

Just go get the package:

go get -u github.com/peterhellberg/jsonstore

Usage

    package main

    import (
    	"context"
    	"fmt"

    	"github.com/peterhellberg/jsonstore"
    )

    const secret = "2ba7860f742fc15d5b6e1508e2de1e0cde2c396f7c52a877905befb4e970eaaf"

    func main() {
    	ctx := context.Background()

    	store := jsonstore.New(jsonstore.Secret(secret))

    	store.Post(ctx, "/example", map[string]interface{}{
    		"number":  1234,
    		"boolean": true,
    		"string":  "example",
    	})

    	store.Put(ctx, "/example/string", "modified")

    	store.Delete(ctx, "/example/boolean")

    	var resp map[string]interface{}

    	store.Get(ctx, "/", &resp)

    	fmt.Printf("%+v\n", resp)
    }

*/
package jsonstore
