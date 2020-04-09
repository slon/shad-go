package exec

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
