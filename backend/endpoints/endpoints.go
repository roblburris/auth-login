package endpoints

import "net/http"

type RequestHandler func(w http.ResponseWriter, r *http.Request)

