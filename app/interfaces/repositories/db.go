package repositories

// DbHandler interface
type DbHandler interface {
	Execute(query string, args ...interface{}) (DbResult, error)
	Query(query string, args ...interface{}) (DbRows, error)
	QueryRow(query string, args ...interface{}) DbRow
	Close() error
}

// DbRow interface
type DbRow interface {
	Scan(dest ...interface{}) error
}

// DbRows interface
type DbRows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

// DbResult summarizes an executed SQL command.
type DbResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
