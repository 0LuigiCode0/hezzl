package dpostgres

import "time"

type Good struct {
	Id          int       `db:"id"`
	ProjectId   int       `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

type InsertGood struct {
	ProjectId int
	Name      string
}

type UpdateGood struct {
	Id          int
	Name        string
	Description string
}

type Meta struct {
	Total   int `db:"total"`
	Removed int `db:"removed"`
}
