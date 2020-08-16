package httpdemo

import (
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
}

func NewHttpServer() (*HttpServer, error) {
	return &HttpServer{}, nil
}

func (hs *HttpServer) Start() error {
	fmt.Println("Starting httpserver")
	http.HandleFunc("/", index)
	fmt.Println("register handleFunc / to index")

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Println("failed to http.ListenAndServe(), err: ", err)
		return err
	} else {
		return nil
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("failed to ParseForm to http.Request, err: ", err)
		return
	}

	fmt.Println("r.Form: ", r.Form)
	fmt.Println("r.URL.Path: ", r.URL.Path)
	fmt.Println("r.URL.Scheme: ", r.URL.Scheme)
	fmt.Println("url_long: ", r.Form.Get("url_long"))
	for k, v := range r.Form {
		fmt.Println("key: ", k, ", value: ", v)
	}
	_, _ = fmt.Fprintf(w, "Thanks for your letter!")
}
