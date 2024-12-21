package authority

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uint64         `json:"id" query:"id" param:"id" form:"id" gorm:"primaryKey,autoIncrement,index,unique,uniqueIndex"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (b Base) IsDeleted() bool {
	return b.DeletedAt.Valid
}

func (b Base) IsValid() bool {
	return b.ID > 0 && !b.CreatedAt.IsZero()
}
