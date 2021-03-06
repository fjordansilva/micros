package go_micro_srv_user

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {

	u := uuid.NewV4()
	/*
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}
	*/
	return scope.SetColumn("Id", u.String())
}
