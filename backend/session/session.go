package session

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
)

type ValueSession struct {
    aud string `json:aud`
    role string `json:role`
    expiration int64 `json:exp`
}

func NewSession(ctx context.Context, red *redis.Client, aud string, role string) (string, error) {
    // creates a new session in the redis DB
    newSessionValue := ValueSession{
        aud: aud,
        role: role,
        expiration: time.Now().AddDate(0, 1, 0).Unix(),
    }
    a := []byte(fmt.Sprintf("%v", newSessionValue))
    byteStruct := new(bytes.Buffer)
    json.NewEncoder(byteStruct).Encode(a)
    return string(a), nil
}

func CheckSession(ctx context.Context, red *redis.Client, key string) {

}