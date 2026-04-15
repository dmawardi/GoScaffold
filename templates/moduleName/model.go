package moduleName

import (
	"time"

	"gorm.io/gorm"
)

type ModuleName struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Add your fields here
	Name string `json:"name,omitempty" binding:"required,min=1,max=255"`
}

type ModuleNameModule struct {
	Repository ModuleNameRepository
	Service    ModuleNameService
	Controller ModuleNameController
}
