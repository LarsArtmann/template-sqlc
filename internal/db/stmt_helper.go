package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// errClosingStatements is a static error for statement close failures.
var errClosingStatements = errors.New("error closing statements")

func CloseStatements(stmts ...*sql.Stmt) error {
	var errs []error

	for _, stmt := range stmts {
		if stmt != nil {
			cerr := stmt.Close()
			if cerr != nil {
				errs = append(errs, cerr)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %d errors occurred", errClosingStatements, len(errs))
	}

	return nil
}
