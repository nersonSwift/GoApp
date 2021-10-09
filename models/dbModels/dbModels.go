package dbModels

type User struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Country string `db:"country"`
}
