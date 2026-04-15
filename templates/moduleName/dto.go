package moduleName

// DTO for creating a new ModuleName
type CreateModuleName struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// DTO for updating a ModuleName
type UpdateModuleName struct {
	Name string `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
}
