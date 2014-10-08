package clevercache

import (
	gctx "code.google.com/p/go.net/context"
	"errors"
	"log"
)

// The client used to interact with the library and the cacheing engine used
type Client struct {
	c Clevercacher // a cacheing library that implements the Clevercacher interface
}

// the key for the cacheing engine stored in the context
var enginekey gctx.Key

type Clevercacher interface {
	Get(ctx gctx.Context, key interface{}) (value interface{}, hit bool, err error)
	Set(ctx gctx.Context, key interface{}, value interface{}) (err error)
	Data(ctx gctx.Context, key interface{}) (value interface{}, err error)
}

func init() {
	enginekey = gctx.NewKey("engine")
}

// Return a new client
func New() *Client {
	return &Client{}
}

// Return a new context with the required cacheing engine embedded
func NewContext(ctx gctx.Context, engine interface{}) gctx.Context {
	return gctx.WithValue(ctx, enginekey, engine)
}

// Get the cacheing engine from the context
func FromContext(ctx gctx.Context) interface{} {
	// ctx.Value returns nil if ctx has no value for the key;
	engine := ctx.Value(enginekey)
	return engine
}

// add the cacheing engine to the client
func (cl *Client) RegisterCache(c Clevercacher) {
	cl.c = c
}

func (cl *Client) Get(ctx gctx.Context, key interface{}) (value interface{}, err error) {

	value, hit, err := cl.c.Get(ctx, key)

	if hit != true {
		// run the cacher in the background
		go func(ctx gctx.Context, cl *Client, key interface{}) {
			value, err := cl.c.Data(ctx, key)
			log.Println("value from data - ", value)
			if err != nil {
				log.Printf("data for key %s cannot be found", key.(string))
			}
			log.Println("setting data")
			err = cl.c.Set(ctx, key, value)

		}(ctx, cl, key)

		return nil, errors.New("no data found - generating new data")
	}

	if err != nil {
		return nil, err
	}

	return value, err

}

func (cl *Client) Set(ctx gctx.Context, key interface{}, value interface{}) (err error) {

	err = cl.c.Set(ctx, key, value)

	return

}
