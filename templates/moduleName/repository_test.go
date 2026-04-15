package moduleName_test

import (
	"fmt"
	"testing"

	"github.com/dmawardi/goTemplate/internal/common"
	"github.com/dmawardi/goTemplate/internal/moduleName"
	"github.com/dmawardi/goTemplate/internal/testutil"
)

func TestModuleNameRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		moduleName    *moduleName.ModuleName
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful moduleName creation",
			moduleName: &moduleName.ModuleName{
				Name: "Test ModuleName",
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "successful moduleName creation with minimal fields",
			moduleName: &moduleName.ModuleName{
				Name: "Minimal ModuleName",
			},
			wantErr:       false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := module.moduleName.repository.Create(tt.moduleName)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("Create() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("Create() error mismatch: got %+v, want %+v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("Create() unexpected error: %v", err)
					return
				}
				// Assert that the created moduleName is not nil and has related details
				testutil.AssertEqual(t, created, tt.moduleName, tt.name)
			}
		})
	}
	// Cleanup: Delete test moduleNames after test completes
	t.Cleanup(func() {
		all, err := module.moduleName.repository.FindAll(10, 0, "", []common.QueryConditionParameters{})
		if err != nil {
			t.Fatalf("Failed to fetch all moduleNames: %v", err)
		}
		var ids []int
		for _, mn := range *all.Data {
			ids = append(ids, int(mn.ID))
		}
		module.moduleName.repository.BulkDelete(ids)
	})
}

func TestModuleNameRepository_FindById(t *testing.T) {
	// Create a test moduleName first
	created, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName: %v", err)
	}

	// Cleanup: Delete test moduleName after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name          string
		id            int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "find existing moduleName",
			id:            int(created.ID),
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:          "find non-existent moduleName",
			id:            99999,
			wantErr:       true,
			expectedError: fmt.Errorf("record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := module.moduleName.repository.FindById(tt.id)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("FindById() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("FindById() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("FindById() unexpected error: %v", err)
					return
				}
				// Assert that the found moduleName is not nil and has related details
				testutil.AssertEqual(t, found, created, tt.name)
			}
		})
	}
}

func TestModuleNameRepository_FindByName(t *testing.T) {
	// Create a test moduleName first
	created, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName: %v", err)
	}

	// Cleanup: Delete test moduleName after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name          string
		search        string
		wantErr       bool
		expectedError error
	}{
		{
			name:          "find existing moduleName",
			search:        created.Name,
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:          "find non-existent moduleName",
			search:        "nonexistent-name",
			wantErr:       true,
			expectedError: fmt.Errorf("record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := module.moduleName.repository.FindByName(tt.search)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("FindByName() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("FindByName() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("FindByName() unexpected error: %v", err)
					return
				}
				// Assert that the found moduleName is not nil and has related details
				testutil.AssertEqual(t, found, created, tt.name)
			}
		})
	}
}

func TestModuleNameRepository_FindAll(t *testing.T) {
	// Create multiple test moduleNames
	moduleName1, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName 1",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName 1: %v", err)
	}

	moduleName2, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName 2",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName 2: %v", err)
	}

	moduleName3, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName 3",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName 3: %v", err)
	}

	// Cleanup: Delete test moduleNames after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(moduleName1.ID), int(moduleName2.ID), int(moduleName3.ID)})
	})

	tests := []struct {
		name          string
		limit         int
		offset        int
		order         string
		conditions    []common.QueryConditionParameters
		wantErr       bool
		expectedError error
		wantCount     int // Expected exact or minimum number of moduleNames
	}{
		{
			name:          "find all moduleNames with limit",
			limit:         2,
			offset:        0,
			order:         "created_at ASC",
			conditions:    []common.QueryConditionParameters{},
			wantErr:       false,
			expectedError: nil,
			wantCount:     2,
		},
		{
			name:          "find all moduleNames with offset",
			limit:         10,
			offset:        1,
			order:         "created_at ASC",
			conditions:    []common.QueryConditionParameters{},
			wantErr:       false,
			expectedError: nil,
			wantCount:     2, // Should have at least 2 remaining after offset
		},
		{
			name:   "find all moduleNames with condition",
			limit:  10,
			offset: 0,
			order:  "created_at ASC",
			conditions: []common.QueryConditionParameters{
				{Condition: "name = ?", Value: "Test ModuleName 3"},
			},
			wantErr:       false,
			expectedError: nil,
			wantCount:     1,
		},
		{
			name:          "find all moduleNames with limit no condition",
			limit:         3,
			offset:        0,
			order:         "created_at ASC",
			conditions:    []common.QueryConditionParameters{},
			wantErr:       false,
			expectedError: nil,
			wantCount:     3, // All test moduleNames with limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := module.moduleName.repository.FindAll(tt.limit, tt.offset, tt.order, tt.conditions)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("FindAll() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("FindAll() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("FindAll() unexpected error: %v", err)
					return
				}

				if result == nil {
					t.Error("FindAll() returned nil result")
					return
				}
				if result.Data == nil {
					t.Error("FindAll() returned nil data")
					return
				}

				actualCount := len(*result.Data)

				// For limited results, check exact count
				if tt.limit > 0 && actualCount != tt.wantCount {
					t.Errorf("FindAll() count mismatch: got %v, want %v", actualCount, tt.wantCount)
				}

				// For offset results, check minimum count
				if tt.offset > 0 && actualCount < tt.wantCount {
					t.Errorf("FindAll() insufficient results after offset: got %v, want at least %v", actualCount, tt.wantCount)
				}

				// Verify metadata is present
				if result.Meta.GetMetaData().Total_Records == 0 && tt.wantCount > 0 {
					t.Error("FindAll() metadata total is 0 when moduleNames exist")
				}
			}
		})
	}
}

func TestModuleNameRepository_Update(t *testing.T) {
	// Create a test moduleName first
	created, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName: %v", err)
	}

	// Cleanup: Delete test moduleName after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name          string
		id            int
		updateData    *moduleName.ModuleName
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful update",
			id:   int(created.ID),
			updateData: &moduleName.ModuleName{
				Name: "Updated Name",
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "update non-existent moduleName",
			id:   99999,
			updateData: &moduleName.ModuleName{
				Name: "Updated Name",
			},
			wantErr:       true,
			expectedError: fmt.Errorf("record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := module.moduleName.repository.Update(tt.id, tt.updateData)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("Update() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("Update() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("Update() unexpected error: %v", err)
					return
				}
				// Assert that the updated moduleName is not nil and has related details
				if updated == nil {
					t.Error("Update() returned nil moduleName")
					return
				}
				if updated.Name != tt.updateData.Name {
					t.Errorf("Update() name not updated: got %v, want %v", updated.Name, tt.updateData.Name)
				}
				if !updated.UpdatedAt.After(updated.CreatedAt) {
					t.Error("Update() UpdatedAt timestamp not properly set")
				}
			}
		})
	}
}

func TestModuleNameRepository_Delete(t *testing.T) {
	// Create a test moduleName first
	created, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName: %v", err)
	}

	// Cleanup: Delete test moduleName after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name          string
		id            int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "successful delete",
			id:            int(created.ID),
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:          "delete non-existent moduleName",
			id:            99999,
			wantErr:       false,
			expectedError: nil, // GORM doesn't return error for deleting non-existent records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := module.moduleName.repository.Delete(tt.id)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("Delete() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("Delete() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("Delete() unexpected error: %v", err)
					return
				}

				// Verify the moduleName is actually deleted (soft delete)
				if tt.id != 99999 {
					_, err := module.moduleName.repository.FindById(tt.id)
					if err == nil {
						t.Error("Delete() moduleName still exists after deletion")
					}
				}
			}
		})
	}
}

func TestModuleNameRepository_BulkDelete(t *testing.T) {
	// Create multiple test moduleNames
	moduleName1, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName 1",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName 1: %v", err)
	}

	moduleName2, err := module.moduleName.repository.Create(&moduleName.ModuleName{
		Name: "Test ModuleName 2",
	})
	if err != nil {
		t.Fatalf("Failed to create test moduleName 2: %v", err)
	}

	// Cleanup: Delete test moduleNames after test completes
	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(moduleName1.ID), int(moduleName2.ID)})
	})

	tests := []struct {
		name          string
		ids           []int
		wantErr       bool
		expectedError error
	}{
		{
			name:          "successful bulk delete",
			ids:           []int{int(moduleName1.ID), int(moduleName2.ID)},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:          "bulk delete non-existent moduleNames",
			ids:           []int{99999, 99998},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:          "bulk delete empty slice",
			ids:           []int{},
			wantErr:       false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := module.moduleName.repository.BulkDelete(tt.ids)

			// If an error is expected, assert that it occurred.
			if tt.wantErr {
				// If no error found
				if err == nil {
					t.Errorf("BulkDelete() expected error but got none")
				} else {
					// If an error is found
					// Check if the error message contains the expected substring
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("BulkDelete() error mismatch: got %v, want %v", err, tt.expectedError)
					}
				}
				return

				// Else If no error is expected, assert that it did not occur.
			} else {
				if err != nil {
					t.Errorf("BulkDelete() unexpected error: %v", err)
					return
				}
			}
		})
	}
}
