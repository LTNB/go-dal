package helper

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type IDatabaseHelper interface {
	GetOne(id interface{}) interface{}
	GetAll() interface{}
	Create(bo interface{}) interface{}
	Update(bo interface{}) interface{}
	Delete(bo interface{}) interface{}
}