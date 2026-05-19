package store

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mike-testut/task-api/internal/models"
)

type UserStore interface {
	CreateUser(username, passwordHash string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type PostgresUserStore struct {
	db *sqlx.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgreUserStore {
	return &PostgresUserStore{
		db: sqlx.NewDb(db, "pgx"),
	}
}

var _ UserStore = (*PostgresUserStore)(nil)

func (s *PostgresUserStore) CreateUser(username, passwordHash string)(models.User, error){
	var user models.User
	query:=`INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username`

	err:= s.db.QueryRowx(query, username, passwordHash).StructScan(&user)
	if err != nil { 
		return models.User{}, fmt.Errorf("could not create user: %v", err)
	
	}
	return user, nil
}

func(s *PostgresUserStore) GetUserByUsername(username string)(models.User,error){
	var user models.User
	query:=`SELECT * FROM users WHERE username = $1`
	err:=s.db.Get(&user,query,username)
	if err != nil {
		if err == sql.ErrNoRows{
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("could not get user: %v", err)
	}
	return user, nil
}