package go_micro_srv_user

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {

	// u, err := uuid.NewV4()
	// if err != nil {
	//	return err
	// }
	u := uuid.Must(uuid.NewV4())
	return scope.SetColumn("Id", u.String())
}
