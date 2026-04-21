package shared

// Queries combines BaseQueries with a direct db field for transaction support.
type Queries struct {
	*BaseQueries

	db DBTX
}

// New creates a new Queries instance with the given database connection.
func New(db DBTX) *Queries {
	return &Queries{
		db:          db,
		BaseQueries: NewBase(db),
	}
}
