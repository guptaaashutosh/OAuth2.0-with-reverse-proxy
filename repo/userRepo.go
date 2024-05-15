package repo

import (

	"github.com/jackc/pgx/v5/pgxpool"
)


type User struct {
	db *pgxpool.Pool
}

func UserRepo(db *pgxpool.Pool) Repositories {
	return &User{
		db: db,
	}
}

func ServiceRepo(db *pgxpool.Pool) Repositories {
	return &User{
		db: db,
	}
}

