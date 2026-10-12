package startup

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func Server() {
	mux := http.NewServeMux()
	ctx, cancelCtx := context.WithCancel(context.Background())

	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listeing for server one: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listeing for server one: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}
