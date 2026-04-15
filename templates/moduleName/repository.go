package moduleName

import (
	"fmt"

	"github.com/dmawardi/goTemplate/internal/common"
	"gorm.io/gorm"
)

type ModuleNameRepository interface {
	// Find a list of all moduleNames in the Database
	FindAll(limit int, offset int, order string, conditions []common.QueryConditionParameters) (*common.BasicPaginatedResponse[ModuleName], error)
	Create(moduleName *ModuleName) (*ModuleName, error)
	Update(int, *ModuleName) (*ModuleName, error)
	Delete(int) error
	BulkDelete([]int) error
	// Find
	FindById(int) (*ModuleName, error)
	FindByName(string) (*ModuleName, error)
}

type moduleNameRepository struct {
	DB *gorm.DB
}

func NewModuleNameRepository(db *gorm.DB) ModuleNameRepository {
	return &moduleNameRepository{db}
}

// Creates a moduleName in the database
func (r *moduleNameRepository) Create(moduleName *ModuleName) (*ModuleName, error) {
	// Create moduleName in database
	result := r.DB.Create(&moduleName)
	if result.Error != nil {
		return nil, fmt.Errorf("failed creating moduleName: %w", result.Error)
	}

	return moduleName, nil
}

// Find a list of moduleNames in the database
func (r *moduleNameRepository) FindAll(limit int, offset int, order string, conditions []common.QueryConditionParameters) (*common.BasicPaginatedResponse[ModuleName], error) {
	// Build meta data for moduleNames
	metaData, err := common.BuildMetaData(r.DB, ModuleName{}, limit, offset, order, conditions)
	if err != nil {
		fmt.Printf("Error building meta data: %s", err)
		return nil, err
	}

	// Find all moduleNames with limit, offset, and order
	var moduleNames []ModuleName
	err = common.QueryAll(r.DB, &moduleNames, limit, offset, order, conditions, []string{})
	if err != nil {
		fmt.Printf("Error finding moduleNames: %s", err)
		return nil, err
	}

	return &common.BasicPaginatedResponse[ModuleName]{
		Data: &moduleNames,
		Meta: *metaData,
	}, nil
}

// Delete a moduleName by ID
func (r *moduleNameRepository) Delete(id int) error {
	// Create an empty ref object of type moduleName
	moduleName := ModuleName{}
	// Check if moduleName exists in db
	result := r.DB.Delete(&moduleName, id)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in deleting moduleName: ", result.Error)
		return result.Error
	}
	// else
	return nil
}

// Bulk delete moduleNames by IDs
func (r *moduleNameRepository) BulkDelete(ids []int) error {
	// Delete moduleNames with specified IDs
	err := common.BulkDeleteByIds(ModuleName{}, ids, r.DB)
	if err != nil {
		fmt.Println("error in deleting moduleNames: ", err)
		return err
	}

	// else
	return nil
}

// Find a moduleName by ID
func (r *moduleNameRepository) FindById(id int) (*ModuleName, error) {
	var moduleName ModuleName
	result := r.DB.First(&moduleName, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &moduleName, nil
}

// Find a moduleName by Name
func (r *moduleNameRepository) FindByName(name string) (*ModuleName, error) {
	var moduleName ModuleName
	result := r.DB.Where("name = ?", name).First(&moduleName)
	if result.Error != nil {
		return nil, result.Error
	}
	return &moduleName, nil
}

// Update a moduleName
func (r *moduleNameRepository) Update(id int, moduleName *ModuleName) (*ModuleName, error) {
	// Find existing moduleName
	existingModuleName, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	result := r.DB.Model(existingModuleName).Updates(moduleName)
	if result.Error != nil {
		return nil, fmt.Errorf("failed updating moduleName: %w", result.Error)
	}

	return existingModuleName, nil
}
