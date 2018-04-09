/*

Package jsonstore is a client for the https://www.jsonstore.io API

Installation

    go get -u github.com/peterhellberg/jsonstore

Usage

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

*/
package jsonstore
