package models

import (
	"database/sql"
	"errors"
	"time"
)

type List struct {
	ID          int
	Title       string
	Description string
	Created     time.Time
	Expires     time.Time
}

type ListModel struct {
	DB *sql.DB
}

func (l *ListModel) Insert(title, description string, expires int) (int, error) {
	stmt := `INSERT INTO lists (title, description, expires)
  VALUES ($1, $2, now() + interval '1 day' * $3) RETURNING id`

	var id int
	err := l.DB.QueryRow(stmt, title, description, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (l *ListModel) Get(id int) (List, error) {
	stmt := `SELECT id, title, description, created, expires FROM lists WHERE id = $1`

	var list List
	err := l.DB.QueryRow(stmt, id).Scan(&list.ID, &list.Title, &list.Description, &list.Created, &list.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return List{}, ErrNoRecord
		} else {
			return List{}, err
		}
	}

	return list, nil
}

func (l *ListModel) Latest() ([]List, error) {
	stmt := `SELECT id, title, description, created, expires FROM lists ORDER BY created DESC LIMIT 5`

	rows, err := l.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

  defer rows.Close()

  var lists []List

  for rows.Next() {
    var l List

    err = rows.Scan(&l.ID, &l.Title, &l.Description, &l.Created, &l.Expires)
    if err != nil {
      return nil, err
    }

    lists = append(lists, l)
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

	return lists, nil
}
