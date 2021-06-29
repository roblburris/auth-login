package db

import (
    "context"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
)

const CHECK_USER_EXISTS = `SELECT u.aud
                        FROM USERS as u
                        WHERE u.email = $1`

const GET_USER = `SELECT u.audd, u.email, u.name, u.role
                        FROM USERS as u, USERS as g
                        WHERE u.aud = g.aud
                        AND u.email = $1
                        AND g.googleID = $2`

const GET_PW_SALT = `SELECT n.salt, n.pw
                    FROM USERS as n
                    WHERE n.aud = $1`

const GET_NON_GISGN_USER = `SELECT u.aud, u.email, u.name, u.role
                        FROM USERS as u, USERS as n
                        WHERE u.aud = n.aud
                        AND u.email = $1
                        AND n.pw = $2`

func CheckUser(ctx context.Context, pool *pgxpool.Pool, email string, aud string) bool {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
        return false
    }
    defer conn.Release()

    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
    })
    if err != nil {
        log.Printf("ERROR: Unable to set transaction level. %v\n", err)
        return false
    }

    defer func(tx pgx.Tx, ctx context.Context) {
        err := tx.Rollback(ctx)
        if err != nil {
            tx.Rollback(ctx)
        }
    }(tx, ctx)
    
    rows, err := tx.Query(ctx, GET_NON_GISGN_USER, "abc", "def")
    // err handling
    if err != nil {
        log.Printf("ERROR: Unable to get blog posts from DB. %v\n", err)
        return false
    }

    defer rows.Close()
    var aud string 
    
    for rows.Next() {
        err = rows.Scan(&aud)
        if err != nil {
            log.Printf("ERROR: Unable to parse SQL data. %v\n", err)
            return false
        }
    }
    if aud == "" {
        tx.Rollback();
    }
    err = tx.Commit(ctx)
    if err != nil {
        tx.Rollback(ctx)
    }
    return true
}

// if non google user
// 1) check to see if email exists in users DB
// 2) if email exists, get hashed_password (byte array) and salt for that user
// 3) hash provided password (string) (func arg) with salt and see if hash equals password
// 4) if equal, user is logged in else return bad

// inside rows.next()
//      var password byte[] 
//      var salt byte[]
//      rows.Scan(&password, &salt)