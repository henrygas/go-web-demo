package httpdemo

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpServer struct {
}

func NewHttpServer() (*HttpServer, error) {
	return &HttpServer{}, nil
}

func (hs *HttpServer) Start() error {
	http.HandleFunc("/", hs.index)
	http.HandleFunc("/login", hs.login)
	http.HandleFunc("/upload", hs.upload)
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
		fmt.Println("r.url_long: ", r.Form["url_long"])
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

func (hs *HttpServer) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles(fmt.Sprintf("%s/login.gtpl", templateDir))
		_ = t.Execute(w, nil)
	} else if err := r.ParseForm(); err != nil {
		fmt.Println("failed to r.ParseForm(), err: ", err)
	} else {
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
		//_, _ = w.Write([]byte(r.Form.Get("username")))
		template.HTMLEscapeString(r.Form.Get("username"))
	}
}

func (hs *HttpServer) upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		_, _ = io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles(fmt.Sprintf("%s/upload.gtpl", templateDir))
		_ = t.Execute(w, token)
	} else if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println("failed to r.ParseMultipartForm, err: ", err)
	} else {
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("failed to r.FormFile(), err: ", err)
			return
		}
		defer file.Close()

		_, err = fmt.Fprintf(w, "%v", handler.Header)
		if err != nil {
			log.Println("failed to fmt.Fprintf(), err: ", err)
			return
		}

		f, err := os.OpenFile(fmt.Sprintf("%s/%s", testDir, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("failed to os.OpenFile(), err: ", err)
			return
		}
		defer f.Close()

		_, _ = io.Copy(f, file)
	}
}
