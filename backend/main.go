package main

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/roblburris/auth-login/endpoints"
    "github.com/go-redis/redis/v8"
    "io/ioutil"
    "log"
    "net/http"
)

const URL = "postgres://localhost:5432/rfts-test"

func main() {
    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, URL)
    if err != nil {
        log.Fatalf("Unable to start pgxpool. Error %v\n", err)
    }
    defer pool.Close()
    conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Fatalf("Unable to connect to DB. Error: %v\n", err)
    }

    deleteTables, err := ioutil.ReadFile("./db/setup/delete-tables.sql")
    if err != nil {
        log.Fatalf("Unable to read delete-tables file, no tests run. Error: %v\n", err)
    }
    _, err = conn.Exec(ctx, string(deleteTables))
    if err != nil {
        log.Fatalf("Unable to wipe existing tables, no tests run. Error: %v\n", err)
    }
    createTables, err := ioutil.ReadFile("./db/setup/create-tables.sql")
    if err != nil {
        log.Fatalf("Unable to read create-tables file, no tests run. Error: %v\n", err)
    }
    _, err = conn.Exec(ctx, string(createTables))
    if err != nil {
        log.Fatalf("Unable create necessary tables, no tests run %v\n", err)
    }

    insertData, err := ioutil.ReadFile("./db/setup/insert-data.sql")
    if err != nil {
        log.Fatalf("Unable to read insert-data file, no tests run. Error: %v\n", err)
    }
    _, err = conn.Exec(ctx, string(insertData))
    if err != nil {
        log.Fatalf("Unable insert data, no tests run %v\n", err)
    }
    conn.Release()

    // setup redis
    red := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	loginHandler := endpoints.LoginEndpoint(ctx, pool, red)
    http.Handle("/", http.FileServer(http.Dir("../frontend/")))
	http.HandleFunc("/login", loginHandler)
	log.Printf("reached\n")
    log.Print(http.ListenAndServe(":8000", nil))
}




