package helper

import (
	"github.com/google/uuid"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

/**
 * base bo for all table, contain id as primary key
 */
type BaseBo struct {
	Id interface{} `json:"id" primary:"id"`
}

func (bo *BaseBo) UUIDGenerate() {
	id := uuid.New()
	bo.Id = id.String()
}

/**
 * versioning data, prevent override data when modify
 */
type Version struct {
	Version string `json:"version"`
}

/**
 * log audit info
 */
type Audit struct {
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}
