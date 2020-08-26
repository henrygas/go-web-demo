package httpdemo

import (
	"context"
	"crypto/md5"
	"fmt"
	"go-web/httpdemo/session"
	_ "go-web/httpdemo/session/provider"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpServer struct {
	ud             *UserDao
	sessionManager *session.Manager
}

func NewHttpServer() (*HttpServer, error) {
	sessionManager, err := session.NewManager("memory", CookieName, 3600)
	if err != nil {
		return nil, err
	}
	go sessionManager.GC()
	return &HttpServer{sessionManager: sessionManager}, nil
}

func (hs *HttpServer) Start() error {
	ud, err := NewUserDao()
	if err != nil {
		log.Println("failed to NewUserDao(), err: ", err)
		return err
	} else {
		hs.ud = ud
	}

	http.HandleFunc("/", hs.index)
	http.HandleFunc("/login", hs.login)
	http.HandleFunc("/register", hs.register)
	http.HandleFunc("/upload", hs.upload)
	if err := http.ListenAndServe(Address, nil); err != nil {
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
	session := hs.sessionManager.SessionStart(w, r)
	if r.Method == "GET" {
		t, _ := template.ParseFiles(fmt.Sprintf("%s/login.gtpl", TemplateDir))
		w.Header().Set("Content-Type", "text/html")
		if session != nil {
			_ = t.Execute(w, session.Get("username"))
		} else {
			_ = t.Execute(w, nil)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("failed to r.ParseForm(), err: ", err)
			return
		}
		if len(r.Form.Get("username")) == 0 {
			_, _ = w.Write([]byte(ERROR_USERNAME_INVALID))
			return
		}
		username := r.Form.Get("username")
		if len(r.Form.Get("password")) == 0 {
			_, _ = w.Write([]byte(ERROR_PASSWORD_INVALID))
			return
		}
		password := r.Form.Get("password")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer func() {
			log.Println("context has been canceled")
			cancel()
		}()

		userinfo, err := hs.ud.GetUserinfoByUsername(ctx, username)
		if err != nil {
			_, _ = w.Write([]byte(ERROR_USERNAME_INVALID))
			return
		}
		if userinfo == nil {
			_, _ = w.Write([]byte(ERROR_USER_NOT_EXIST))
			return
		}
		if userinfo.Password != password {
			_, _ = w.Write([]byte(ERROR_PASSWORD_INVALID))
			return
		}

		_ = session.Set("username", username)
		http.Redirect(w, r, "/upload", 302)
	}
}

func (hs *HttpServer) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(fmt.Sprintf("%s/register.gtpl", TemplateDir))
		err := t.Execute(w, nil)
		if err != nil {
			log.Println("failed to t.Execute, err: ", err)
			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println("failed to r.ParseForm(), err: ", err)
			return
		}
		if len(r.Form.Get("username")) == 0 {
			_, err := w.Write([]byte(ERROR_USERNAME_INVALID))
			if err != nil {
				log.Println("failed to w.Write, err: ", err)
				return
			}
		}
		username := r.Form.Get("username")

		if len(r.Form.Get("password")) == 0 {
			_, err = w.Write([]byte(ERROR_PASSWORD_INVALID))
			if err != nil {
				log.Println("failed to w.Write, err: ", err)
				return
			}
		}
		password := r.Form.Get("password")

		err = hs.ud.CreateUserinfo(context.Background(), &Userinfo{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Printf("failed to CreateUserInfo(username=%s, password=%s), err: %s\n",
				username, password, err.Error())
			return
		}

		http.Redirect(w, r, "/login", 302)
	}
}

func (hs *HttpServer) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		_, _ = io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles(fmt.Sprintf("%s/upload.gtpl", TemplateDir))
		_ = t.Execute(w, token)
	} else {
		session, err := hs.sessionManager.GetSession(r)
		if err != nil || session == nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Println("failed to r.ParseMultipartForm, err: ", err)
		} else {
			file, fh, err := r.FormFile("uploadfile")
			if err != nil {
				fmt.Println("failed to r.FormFile(), err: ", err)
				return
			}
			defer file.Close()

			_, err = fmt.Fprintf(w, "%v", fh.Header)
			if err != nil {
				log.Println("failed to fmt.Fprintf(), err: ", err)
				return
			}

			f, err := os.OpenFile(fmt.Sprintf("%s/%s", TestDir, fh.Filename), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println("failed to os.OpenFile(), err: ", err)
				return
			}
			defer f.Close()

			_, _ = io.Copy(f, file)
		}
	}
}
