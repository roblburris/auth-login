package db

import (
    "context"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
)

const CHECK_USER_EXISTS = `SELECT u.uid
                        FROM USERS as u
                        WHERE u.email = $1`

const GET_GSUITE_USER = `SELECT u.uid, u.email, u.name, u.role
                        FROM USERS as u, GSUITE_USERS as g
                        WHERE u.uid = g.uid
                        AND u.email = $1
                        AND g.googleID = $2`

const GET_PW_SALT = `SELECT n.salt
                    FROM NON_GSUITE_USERS as n
                    WHERE n.uid = $1`

const GET_NON_GSUITE_USER = `SELECT u.uid, u.email, u.name, u.role
                        FROM USERS as u, NON_GSUITE_USERS as n
                        WHERE u.uid = n.uid
                        AND u.email = $1
                        AND n.pw = $2`

func CheckGsuiteUser(ctx context.Context, pool *pgxpool.Pool, email string, googleID int) {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
        return
    }
    defer conn.Release()

    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
    })
    if err != nil {
        log.Printf("ERROR: Unable to set transaction level. %v\n", err)
        return
    }

    defer func(tx pgx.Tx, ctx context.Context) {
        err := tx.Rollback(ctx)
        if err != nil {
            tx.Rollback(ctx)
        }
    }(tx, ctx)


}