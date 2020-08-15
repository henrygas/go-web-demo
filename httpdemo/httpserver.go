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
	http.HandleFunc("/", hs.index)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Println("failed to http.ListenAndServe, err: ", err)
		return err
	} else {
		return nil
	}
}

func (hs *HttpServer) index(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("failed to parseForm, err: ", err)
		return
	} else {
		fmt.Println("r.Form: ", r.Form)
		fmt.Println("r.Path: ", r.URL.Path)
		fmt.Println("r.Scheme: ", r.URL.Scheme)
		fmt.Println("r.url_long: ", r.Form.Get("url_long"))
		for k, v := range r.Form {
			fmt.Println("key: ", k, ", value: ", v)
		}
		if n, err := fmt.Fprintf(w, "Hello from henry server"); err != nil {
			fmt.Println("failed to write content to response, err: ", err)
			return
		} else {
			fmt.Printf("henry server has written %d bytes to client\n", n)
		}
	}
}
