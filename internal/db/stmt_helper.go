package db

import (
	"database/sql"
	"fmt"
)

func CloseStatements(stmts ...*sql.Stmt) error {
	var err error
	for _, stmt := range stmts {
		if stmt != nil {
			if cerr := stmt.Close(); cerr != nil {
				err = fmt.Errorf("error closing statement: %w", cerr)
			}
		}
	}
	return err
}
