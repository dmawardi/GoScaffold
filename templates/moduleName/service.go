package moduleName

import (
	"fmt"

	"github.com/dmawardi/goTemplate/internal/common"
)

type ModuleNameService interface {
	// CRUD operations
	FindAll(limit int, offset int, order string, conditions []common.QueryConditionParameters) (*common.BasicPaginatedResponse[ModuleName], error)
	FindById(int) (*ModuleName, error)
	FindByName(string) (*ModuleName, error)
	Create(moduleName *CreateModuleName) (*ModuleName, error)
	Update(int, *UpdateModuleName) (*ModuleName, error)
	Delete(int) error
	BulkDelete([]int) error
}

type moduleNameService struct {
	repo ModuleNameRepository
}

// Builds a new service with injected repository
func NewModuleNameService(repo ModuleNameRepository) ModuleNameService {
	return &moduleNameService{repo: repo}
}

// Creates a moduleName in the database
func (s *moduleNameService) Create(moduleName *CreateModuleName) (*ModuleName, error) {
	// Create a new moduleName of type db ModuleName
	toCreate := ModuleName{
		Name: moduleName.Name,
	}

	// Create above moduleName in database
	created, err := s.repo.Create(&toCreate)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// Find all moduleNames
func (s *moduleNameService) FindAll(limit int, offset int, order string, conditions []common.QueryConditionParameters) (*common.BasicPaginatedResponse[ModuleName], error) {
	result, err := s.repo.FindAll(limit, offset, order, conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to find moduleNames: %w", err)
	}

	return result, nil
}

// Find a moduleName by ID
func (s *moduleNameService) FindById(id int) (*ModuleName, error) {
	moduleName, err := s.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find moduleName by id: %w", err)
	}

	return moduleName, nil
}

// Find a moduleName by Name
func (s *moduleNameService) FindByName(name string) (*ModuleName, error) {
	moduleName, err := s.repo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to find moduleName by name: %w", err)
	}

	return moduleName, nil
}

// Update a moduleName in the database
func (s *moduleNameService) Update(id int, moduleName *UpdateModuleName) (*ModuleName, error) {
	// Create db ModuleName type from incoming DTO
	toUpdate := &ModuleName{
		Name: moduleName.Name,
	}

	// Update using repo
	updated, err := s.repo.Update(id, toUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update moduleName: %w", err)
	}

	return updated, nil
}

// Delete moduleName in database
func (s *moduleNameService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete moduleName: %w", err)
	}

	return nil
}

// Deletes multiple moduleNames in database
func (s *moduleNameService) BulkDelete(ids []int) error {
	err := s.repo.BulkDelete(ids)
	if err != nil {
		return fmt.Errorf("failed to bulk delete moduleNames: %w", err)
	}

	return nil
}
