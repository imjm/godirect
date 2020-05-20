package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8000", "HTTP service address")
		url      = flag.String("url", "", "Redirect URL")
	)
	flag.Parse()

	if *url == "" {
		log.Fatal("Redirect url is empty")
	}

	h := handler{
		URL: *url,
	}

	serve(*httpAddr, h)
}

func serve(addr string, h handler) {
	fasthttp.ListenAndServe(addr, h.HandleHTTP)
}

type handler struct {
	URL string
}

func (h *handler) HandleHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Println(string(ctx.Path()))
	statusCode, body, err := fasthttp.Get(nil, h.URL+string(ctx.Path()))
	if err != nil {
		return
	}
	ctx.SetStatusCode(statusCode)
	fmt.Fprintf(ctx, string(body))
}
