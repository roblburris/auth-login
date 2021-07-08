package db

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

const INSERT_USER = `INSERT INTO USERS(aud, email, password, salt, role)
                        VALUES ( $1,$2, $3, $4, $5 )`

const CHECK_AUD = `SELECT COUNT(*)
					FROM USERS as u
					WHERE u.aud = $1`

func checkAudExists(ctx context.Context, pool *pgxpool.Pool, aud string) (bool, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
		return false, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadUncommitted,
	})
	if err != nil {
		log.Printf("ERROR: Unable to set transaction level. %v\n", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, CHECK_AUD, aud)
	if err != nil {
		log.Printf("ERROR: Unable to query Users DB to check aud. %v\n", err)
		return false, err
	}

	var count uint32
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Printf("ERROR: unable to query Users DB to check aud. %v\n", err)
			return false, err
		}
	}
	return count != 0, nil
}

func InsertUser(ctx context.Context, pool *pgxpool.Pool, aud string, email string, password []byte, role string) (bool, error) {
    // first check if user exists in DB
	testAud, err := getUserAud(ctx, pool, email)
	if testAud != "" {
		return false, errors.New("email_exists")
	}

	conn, err := pool.Acquire(ctx)
    if err != nil {
        log.Printf("ERROR: unable to acquire connection from pool. %v\n", err)
        return false, err
    }
    defer conn.Release()

    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
    })
    if err != nil {
        log.Printf("ERROR: Unable to set transaction level. %v\n", err)
        return false, err
    }
    defer tx.Rollback(ctx)

	if aud == "" {
		salt := generateRandomSalt(32)

		h := sha1.New()
		h.Write(append(password, salt...))
		hashedPwd := h.Sum(nil)

		h = sha1.New()
		h.Write(append([]byte(email), salt...))
		aud = string(h.Sum(nil))

		// ensure aud is unique
		audExists, err := checkAudExists(ctx, pool, aud)
		if err != nil {
			log.Printf("ERROR: unable to check whether generated aud exists.\n")
			return false, err
		}
		for audExists {
			h = sha1.New()
			salt = generateRandomSalt(32)
			h.Write(append([]byte(aud), salt...))
			aud = string(h.Sum(nil))
			audExists, err = checkAudExists(ctx, pool, aud)
			if err != nil {
				log.Printf("ERROR: unable to check whether generated aud exists.\n")
				return false, err
			}
		}

		_, err = tx.Query(ctx, INSERT_USER, aud, email, hashedPwd, salt, role)
	
		if err != nil {
			log.Printf("ERROR: Unable to query Users DB to get aud. %v\n", err)
			return false, err
		}
	} else {
		_, err := tx.Query(ctx, INSERT_USER, aud, email, nil, nil, role)
		if err != nil {
			log.Printf("ERROR: Unable to insert User info into DB. %v\n", err)
			return false , err
		}
	}

	_ = tx.Commit(ctx)
	return true, nil
}

func generateRandomSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	
	if err != nil {
		log.Printf("ERROR: unable to generate sakt")
	}
	
	return salt
}