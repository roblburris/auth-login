package db

import (
    "bytes"
    "context"
    "crypto/sha1"
    "errors"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
)

const GET_USER_AUD = `SELECT u.aud
                        FROM USERS as u
                        WHERE u.email = $1`

const GET_USER = `SELECT u.aud, u.email, u.name, u.role
                        FROM USERS as u, USERS as g
                        WHERE u.aud = g.aud
                        AND u.email = $1
                        AND g.googleID = $2`

const GET_PW_SALT = `SELECT n.salt, n.pw
                    FROM USERS as n
                    WHERE n.aud = $1`

const GET_NON_GSIGN_USER = `SELECT u.aud, u.email, u.name, u.role
                        FROM USERS as u, USERS as n
                        WHERE u.aud = n.aud
                        AND u.email = $1
                        AND n.pw = $2`

func CheckGsuiteUser(ctx context.Context, pool *pgxpool.Pool,
    email string, googleAud string) (bool, error) {
    // first check to see if user exists in DB
    aud, err := getUserAud(ctx, pool, email)
    if err != nil {
        log.Printf("ERROR: unable to get aud.\n")
        return false, err
    }

    if aud == "" {
        log.Printf("aud does not exist in Users DB.\n")
        return false, errors.New("USER_DNE")
    }

    // we have aud, check to see if it matches googleAUD
    return aud == googleAud, nil
}


func CheckNonGsuiteUser(ctx context.Context, pool *pgxpool.Pool,
    email string, pwd string) (bool, error) {
    // first check to see if the user exists in the DB
    aud, err := getUserAud(ctx, pool, email)
    if err != nil {
        log.Printf("ERROR: unable to get aud.\n")
        return false, err
    }

    if aud == "" {
        log.Printf("aud does not exist in Users DB.\n")
        return false, errors.New("USER_DNE")
    }

    // we have the aud, get salt and hashedPwd for that user
    salt, hashedPwd, err := getPwd(ctx, pool, aud)
    if err != nil {
        log.Printf("Unable to get salt and hashedPW for given user.")
        return false, err
    }

    if salt == nil || hashedPwd == nil {
        log.Printf("User does not have a stored password. Maybe they signed up using Gsuite?")
        return false, errors.New("NO_PW")
    }

    // hash given password with salt and return whether they are equal
    h := sha1.New()
    h.Write(append([]byte(pwd), salt...))
    return bytes.Equal(h.Sum(nil), hashedPwd), nil
}

func getUserAud(ctx context.Context, pool *pgxpool.Pool, email string) (string, error) {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
        return "", err
    }
    defer conn.Release()

    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
    })
    if err != nil {
        log.Printf("ERROR: Unable to set transaction level. %v\n", err)
        return "", err
    }
    defer tx.Rollback(ctx)

    rows, err := tx.Query(ctx, GET_USER_AUD, email)
    if err != nil {
        log.Printf("ERROR: Unable to query Users DB to get aud. %v\n", err)
        return "", err
    }

    var aud string
    for rows.Next() {
        err = rows.Scan(&aud)
        if err != nil {
            log.Printf("ERROR: unable to unpack results of query for getting aud. %v\n")
            return "", err
        }
    }

    err = tx.Commit(ctx)
    if err != nil {
        log.Printf("ERROR: unable to commit results. %v\n")
        _ = tx.Commit(ctx)
    }
    return aud, nil
}

func getPwd(ctx context.Context, pool *pgxpool.Pool, aud string) ([]byte, []byte, error) {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
        return nil, nil, err
    }

    defer conn.Release()

    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
    })
    if err != nil {
        log.Printf("ERROR: Unable to set transaction level. %v\n", err)
        return nil, nil, err
    }

    defer tx.Rollback(ctx)

    rows, err := tx.Query(ctx, GET_PW_SALT, aud)
    // err handling
    if err != nil {
        log.Printf("ERROR: Unable to get blog posts from DB. %v\n", err)
        return nil, nil, err
    }

    defer rows.Close()
    var salt []byte
    var pwd []byte
    for rows.Next() {
      err = rows.Scan(&salt, &pwd)
      if err != nil {
          log.Printf("ERROR: Unable to parse SQL data. %v\n", err)
          return nil, nil, err
      }
    }

    err = tx.Commit(ctx)
    if err != nil {
      _ = tx.Rollback(ctx)
    }

    return salt, pwd, nil
}