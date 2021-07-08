package session

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
    "log"
    "crypto/sha256"
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
    
    encodedSessionValue, err := json.Marshal(newSessionValue)
    if err != nil {
        log.Printf("ERROR: unable to encode struct. %v\n", err)
        return "", err
    }
    
    
    hash := sha256.New()
    hash.Write([]byte(newSessionValue.aud + string(newSessionValue.expiration)))
    key := string(hash.Sum(nil))

    err = red.Set(ctx, key, encodedSessionValue, 0).Err()
    if err != nil {
        log.Printf("ERROR: error inserting into RedisDB. %v\n", err)
        return "", err
    }
    
    return key, nil
}

func CheckSession(ctx context.Context, red *redis.Client, key string) {

}