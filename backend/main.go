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
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("../frontend/")))
    log.Fatal(http.ListenAndServe(":63342", nil))
}


func (cfg config) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// parse the GoogleJWT that was POSTed from the front-end
	type parameters struct {
		GoogleJWT *string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// Validate the JWT is valid
	claims, err := ValidateGoogleJWT(*params.GoogleJWT)
	if err != nil {
		respondWithError(w, 403, "Invalid google auth")
		return
	}
	if claims.Email != user.Email {
		respondWithError(w, 403, "Emails don't match")
		return
	}

	// create a JWT for OUR app and give it back to the client for future requests
	tokenString, err := MakeJWT(claims.Email, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, 500, "Couldn't make authentication token")
		return
	}

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})
}

