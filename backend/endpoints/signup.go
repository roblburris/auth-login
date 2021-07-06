package endpoints

// import (
//     "encoding/json"
//     "fmt"
//     "github.com/roblburris/auth-login/auth"
//     "io/ioutil"
//     "log"
//     "net/http"
// )

// func SignupEndpoint() RequestHandler {
// 	return func(w http.ResponseWriter, r *http.Request) {
//         if r.Method != http.MethodPost {
//             text := "405: expected Post"
//             log.Printf("incorrect request, %s\n", text)
//             w.WriteHeader(http.StatusMethodNotAllowed)
//             return
//         }

//         var body map[string]interface{}
//         temp, err := ioutil.ReadAll(r.Body)
//         if err != nil {
//             log.Printf("Unable to decode body. Bad Request. Error: %v\n", err)
//             w.WriteHeader(http.StatusBadRequest)
//             return
//         }
//         err = json.Unmarshal(temp, &body)
//         if err != nil {
//             log.Printf("Unable to decode body. Bad Request. Error: %v\n", err)
//             w.WriteHeader(http.StatusBadRequest)
//             return
//         }

//         if body["type"] == nil {
//             log.Printf("Request not formatted correctly.")
//             w.WriteHeader(http.StatusBadRequest)
//             return
//         }

//         typeOf := fmt.Sprintf("%v", body["type"])
//         if !(typeOf == "gsuite" || typeOf == "non_gsuite") {
//             log.Printf("Request not formatted correctly. Expecting `gsuite` or `non-gsuite` in body[`type`] but got `%s`\n",
//             typeOf )
//             w.WriteHeader(http.StatusBadRequest)
//             return
//         }

//         // case if gsuite user
//         if typeOf  == "gsuite" {
//             payload, err := auth.ValidateGoogleJWT(fmt.Sprintf("%v", body["jwt"]))
//             if err != nil {
//                 log.Printf("does not work. %v\n", err)
//                 w.WriteHeader(http.StatusInternalServerError)
//                 return
//             }
//             log.Printf("%v\n", payload)
//             log.Printf("aud: %v\n", payload.Audience)
            

//         } else { // case if normal user
//             w.WriteHeader(http.StatusInternalServerError)
//             return
//         }

//         // set cookie and return

//         w.WriteHeader(http.StatusOK)
//         return
//     }
// }