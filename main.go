package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DB_USER     = "redhwan"
	DB_PASSWORD = "Pass"
	DB_NAME     = "first_test"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string // Raw text of unparsed attribute-value pairs
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	provifer    Provider
	maxlifetime int64
}


type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}


type MyMux struct{}





func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	secs := now.Unix()

	fmt.Printf("time: %s path: %v scheme: %s\n", time.Unix(secs, 0), r.URL.Path, r.URL.Scheme)
	switch r.URL.Path {
	case "/":
		sayHelloName(w, r)
		return
	case "/login":
		login(w, r)
		return
	case "/register":
		register(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}

}

func NewSessionManager(providerName, cookieName string, maxlifetime int64) (*SessionManager, error) {
	provider, ok := providers[providerName]

	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q (forgotten import?)", providerName)
	}
	return &SessionManager{provifer, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method) // get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("gtpl/home.gtpl")
		err := t.Execute(w, nil)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	} else {
		template.HTMLEscape(w, []byte("Unauthorized login")) // respond to client
	}

	fmt.Fprint(w, "Hello astaxie.!")

}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method) // get request method
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token: ", token)

		expiration := time.Now().Add(10 * 24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
		http.SetCookie(w, &cookie)

		t, _ := template.ParseFiles("gtpl/login.gtpl")
		err := t.Execute(w, token)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	} else {
		// log in request
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			// check token validity
			// fileUpload(w, r)
			template.HTMLEscape(w, []byte("Welcome Back "+r.Form.Get("username"))) // respond to client
			for _, cookie := range r.Cookies() {
				fmt.Fprint(w, cookie.Name)
			}

		} else {
			// give error if no token
			template.HTMLEscape(w, []byte("Unauthorized login")) // respond to client
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print in server side
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method) // get request method
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token: ", token)

		t, _ := template.ParseFiles("gtpl/register.gtpl")
		err := t.Execute(w, token)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	} else {
		// log in request
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			// check token validity
			// fileUpload(w, r)
			template.HTMLEscape(w, []byte("Sucessfully Registered"+r.Form.Get("username"))) // respond to client
		} else {
			// give error if no token
			template.HTMLEscape(w, []byte("Unauthorized login")) // respond to client
		}
		fmt.Println("New Register:")
		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("vla:", strings.Join(v, ""))
		}
	}
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("gtpl/fileUpload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func dbCreateTable(tblName string) {
	var strquery string = ""

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)

	defer db.Close()

	switch tblName {
	case "sessions":
		strquery = "CREATE TABLE sessions (" +
			"sessid varchar(20) CONSTRAINT firstkey PRIMARY KEY," +
			"user_id char(20)," +
			"ip CIDR, timestart timestamp)"
	}

	stmt, err := db.Prepare(strquery)

	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	fmt.Println(res)
}

func main() {
	mux := &MyMux{} // custom router
	// http.HandlerFunc("/", sayHelloName)      //set Route
	// http.HandleFunc("/login", login)
	// http.HandleFunc("/upload", fileUpload)
	err := http.ListenAndServe(":9090", mux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


	//creating global session manager
	var globalSessions *session.SessionManager
    // Then, initialize the session manager
    func init() {
        globalSessions = NewSessionManager("memory","gosessionid",3600)
    }
}
