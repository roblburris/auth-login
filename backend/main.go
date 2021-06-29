// package main

// import (
//     "log"
//     "net/http"
//     "github.com/gorilla/mux"
// )

// func main() {
	// r := mux.NewRouter()
//     r.http.Handle("/", http.FileServer(http.Dir("../frontend/")))
//     err := http.ListenAndServe(":63342", nil)
//     if err != nil {
//         log.Printf("err: %v\n", err)
//     }
//     log.Fatalln("listening")
// }

package main

import (
	"github.com/roblburris/auth-login/endpoints"
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	loginHandler := endpoints.LoginEndpoint()
    http.Handle("/", http.FileServer(http.Dir("../frontend/")))
	http.HandleFunc("/login", loginHandler)
    log.Fatal(http.ListenAndServe(":63342", nil))
}




