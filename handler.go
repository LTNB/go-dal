package go_dal

import "database/sql"

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
type IDatabaseHelper interface {
	GetDatabase() *sql.DB
}