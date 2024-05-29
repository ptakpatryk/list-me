package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
  ID int
  Name string
  Email string
  HashedPassword []byte
  Created time.Time
}

type UserModel struct {
  DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error  {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
  if err != nil {
    return err
  }

  stmt := `INSERT INTO users (name, email, hashed_password) 
  VALUES ($1, $2, $3)`

  
  _, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
  if err != nil {
    var psqlError *pgconn.PgError
    if errors.As(err, &psqlError) {
      if psqlError.Code == "23505" {
        return ErrDuplicateEmail
      }
    }

    return err
  }
  return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error)  {
  return 0, nil
}


func (m *UserModel) Exists(id int) (bool, error)  {
  return false, nil
}
