package endpoints

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func LoginEndpoint() RequestHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            text := "405: expected Post"
            log.Printf("incorrect request, %s\n", text)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        var body map[string]interface{}
        temp, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Printf("Unable to decode body. Bad Request. Error: %v\n", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        err = json.Unmarshal(temp, &body)
        if err != nil {
            log.Printf("Unable to decode body. Bad Request. Error: %v\n", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if body["type"] == nil {
            log.Printf("Request not formatted correctly.")
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if fmt.Sprintf("%v", body["type"]) != "gsuite" ||  fmt.Sprintf("%v", body["type"]) != "non-gsuite" {
            log.Printf("Request not formatted correctly. Expecting `gsuite` or `non-gsuite` in body[`type`] but got %s\n",
                fmt.Sprintf("%v", body["type"]))
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // case if gsuite user
        if fmt.Sprintf("%v", body["type"]) == "gsuite" {

        } else { // case if normal user

        }

        // set cookie and return


        w.WriteHeader(http.StatusOK)
        return
    }
}