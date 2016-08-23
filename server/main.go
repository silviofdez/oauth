package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"gopkg.in/session.v1"
)

var (
	globalSessions *session.Manager
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Location string `json:"location"`
}

func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()
}

func main() {
	manager := manage.NewDefaultManager()

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// token configuration
	cfg := &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
		// refresh token expiration time
		RefreshTokenExp: time.Hour * 24 * 3,
		// whether to generate the refreshing token
		IsGenerateRefresh: true,
	}

	manager.SetAuthorizeCodeTokenCfg(cfg)

	// client store
	manager.MapClientStorage(store.NewTestClientStore(&models.Client{
		ID:     "bawdy-reindeers-14-56dd2bcc2ba94",
		Secret: "6454acedc7024fdfa743c5407da7ad44",
		Domain: "http://localhost:3000",
	}))

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(func(err error) {
		fmt.Println("internal error:", err.Error())
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		accessToken := r.Form.Get("access_token")
		token, err := manager.LoadAccessToken(accessToken)

		if err == nil {

			user := &User{
				token.GetUserID(),
				"test@test.com",
				"localhost",
			}
			jData, _ := json.Marshal(user)
			w.Write(jData)

		} else {
			w.WriteHeader(http.StatusForbidden)
		}

	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("\n\n%#v\n\n", r)
		//
		// body, err := ioutil.ReadAll(r.Body)
		//
		// if err != nil {
		// 	panic(err.Error())
		// }
		//
		// fmt.Printf("\n\n%#v\n\n", string(body))
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			fmt.Printf("\n\n\n%#v\n\n\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	us, err := globalSessions.SessionStart(w, r)

	uid := us.Get("UserID")
	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		us.Set("Form", r.Form)
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	// us.Delete("UserID")
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		us, err := globalSessions.SessionStart(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		us.Set("UserID", "000000")
		us.Set("email", "test@test.com")
		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	us, err := globalSessions.SessionStart(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if us.Get("UserID") == nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	if r.Method == "POST" {
		form := us.Get("Form").(url.Values)
		u := new(url.URL)
		u.Path = "/authorize"
		u.RawQuery = form.Encode()
		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusFound)
		us.Delete("Form")
		return
	}
	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
