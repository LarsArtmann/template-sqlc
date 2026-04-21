package shared

type Queries struct {
	*BaseQueries
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{
		db:          db,
		BaseQueries: NewBase(db),
	}
}
