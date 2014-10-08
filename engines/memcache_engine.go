// This is an example cacheing engine for use with the clevercache package
package main

import (
	gctx "code.google.com/p/go.net/context"
	"github.com/bradfitz/gomemcache/memcache"
)

type mcengine struct {
	mc *memcache.Client
}

func (mc *mcengine) Get(ctx gctx.Context, key interface{}) (value interface{}, hit bool, err error) {

	out, err := mc.mc.Get(key.(string))
	hit = true

	if err != nil {
		hit = false
		return nil, hit, err
	}

	value = out.Value

	return

}

func (mc *mcengine) Set(ctx gctx.Context, key interface{}, value interface{}, params map[string]string) (err error) {

	err = mc.mc.Set(&memcache.Item{Key: key.(string), Value: []byte(value.(string))})

	return
}

func (mc *mcengine) Data(ctx gctx.Context, key interface{}, params map[string]string) (value interface{}, err error) {

	// At this point you should run off and get some data based on the key
	// if this isn't possible you should add some data to the context to allow you to get whatever data you need from what ever 3rd party source

	value = "test string"
	err = nil

	return

}
