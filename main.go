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
		header   = flag.String("header", "", "Redirect Header key")
	)
	flag.Parse()

	if *url == "" {
		log.Fatal("Redirect url is empty")
	}

	h := handler{
		URL:    *url,
		Header: *header,
	}

	serve(*httpAddr, h)
}

func serve(addr string, h handler) {
	fasthttp.ListenAndServe(addr, h.HandleHTTP)
}

type handler struct {
	URL    string
	Header string
}

func (h *handler) HandleHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Println(string(ctx.Path()))

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(h.URL + string(ctx.Path()))

	if h.Header != "" {
		req.Header.SetBytesV(h.Header, ctx.Request.Header.Peek(h.Header))
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		fmt.Errorf("err: %s", err.Error())
	}
	body := resp.Body()
	fmt.Fprintf(ctx, string(body))
}
