package postgres

import (
	"database/sql"
	go_dal "github.com/LTNB/go-dal"
	_ "github.com/lib/pq"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type PostgresHelper struct {
}

func (postgresHelper PostgresHelper) GetDatabase() *sql.DB {
	return go_dal.GetDatabase()
}