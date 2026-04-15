package moduleName_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dmawardi/goTemplate/internal/moduleName"
	"github.com/dmawardi/goTemplate/internal/testutil"
	"github.com/gin-gonic/gin"
)

func TestModuleNameController_FindAll(t *testing.T) {
	// Set up the test environment
	created1, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 1",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}
	created2, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 2",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}
	created3, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 3",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created1.ID), int(created2.ID), int(created3.ID)})
	})

	// Create a test HTTP request and response recorder
	c, w := testutil.CreateTestContext("GET", "/moduleNames", nil)

	// Call the controller method directly
	module.moduleName.controller.FindAll(c)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var responseData map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&responseData); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response contains data field
	data, ok := responseData["data"]
	if !ok {
		t.Fatal("response missing 'data' field")
	}

	// Convert response data back to moduleName structs for assertion
	var responseModuleNames []moduleName.ModuleName
	moduleNamesJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal moduleNames data: %v", err)
	}
	if err := json.Unmarshal(moduleNamesJSON, &responseModuleNames); err != nil {
		t.Fatalf("failed to unmarshal to moduleName structs: %v", err)
	}
	if len(responseModuleNames) < 3 {
		t.Errorf("Expected at least 3 moduleNames in response, got %d", len(responseModuleNames))
	}

	expectedModuleNameDetails := map[uint]*moduleName.ModuleName{
		created1.ID: created1,
		created2.ID: created2,
		created3.ID: created3,
	}

	// Iterate through response moduleNames and compare with expected moduleNames by ID
	for _, moduleNameBeingSearched := range responseModuleNames {
		// If the moduleName ID from the response matches one of the created moduleNames, compare their fields
		if expectedModuleNameDetails[moduleNameBeingSearched.ID] != nil {
			testutil.AssertFieldsEqual(t, moduleNameBeingSearched, *expectedModuleNameDetails[moduleNameBeingSearched.ID], "FindAll")
		} else {
			t.Errorf("Unexpected moduleName ID in response: %d", moduleNameBeingSearched.ID)
		}
	}

	// Verify pagination metadata if present
	if pagination, ok := responseData["pagination"]; ok {
		paginationMap, _ := pagination.(map[string]interface{})
		if totalRecords, ok := paginationMap["total_records"].(float64); ok {
			if totalRecords < 3 {
				t.Errorf("pagination total_records = %v; want >= 3", totalRecords)
			}
		}
	}
}

func TestModuleNameController_Find(t *testing.T) {
	// Create a test moduleName
	created, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Find Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name           string
		id             uint
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:           "find existing moduleName",
			id:             created.ID,
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "find non-existent moduleName",
			id:             99999,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := testutil.CreateTestContext("GET", fmt.Sprintf("/moduleNames/%d", tt.id), nil)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", tt.id)})

			module.moduleName.controller.Find(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %d; want %d", w.Code, tt.expectedStatus)
			}

			if !tt.shouldFail {
				var responseModuleName moduleName.ModuleName
				if err := json.NewDecoder(w.Body).Decode(&responseModuleName); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				testutil.AssertFieldsEqual(t, responseModuleName, *created, tt.name)
			}
		})
	}
}

func TestModuleNameController_Create(t *testing.T) {
	tests := []struct {
		name           string
		createData     *moduleName.CreateModuleName
		expectedStatus int
		shouldFail     bool
	}{
		{
			name: "successful moduleName creation",
			createData: &moduleName.CreateModuleName{
				Name: "New ModuleName",
			},
			expectedStatus: http.StatusCreated,
			shouldFail:     false,
		},
		{
			name: "create moduleName with empty name",
			createData: &moduleName.CreateModuleName{
				Name: "",
			},
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.createData)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			c, w := testutil.CreateTestContext("POST", "/moduleNames", strings.NewReader(string(bodyBytes)))

			module.moduleName.controller.Create(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %d; want %d", w.Code, tt.expectedStatus)
			}

			if !tt.shouldFail {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName creation successful!" {
					t.Error("expected success message in response, instead got:", response)
				}
			}

			t.Cleanup(func() {
				foundModuleName, err := module.moduleName.service.FindByName(tt.createData.Name)
				if err == nil {
					module.moduleName.repository.BulkDelete([]int{int(foundModuleName.ID)})
				}
			})
		})
	}
}

func TestModuleNameController_Update(t *testing.T) {
	// Create a test moduleName
	created, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Update Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	t.Cleanup(func() {
		module.moduleName.repository.BulkDelete([]int{int(created.ID)})
	})

	tests := []struct {
		name           string
		id             uint
		updateData     *moduleName.UpdateModuleName
		expectedStatus int
		shouldFail     bool
	}{
		{
			name: "successful update",
			id:   created.ID,
			updateData: &moduleName.UpdateModuleName{
				Name: "Updated ModuleName",
			},
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name: "update non-existent moduleName",
			id:   99999,
			updateData: &moduleName.UpdateModuleName{
				Name: "Updated ModuleName",
			},
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.updateData)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			c, w := testutil.CreateTestContext("PUT", fmt.Sprintf("/moduleNames/%d", tt.id), strings.NewReader(string(bodyBytes)))
			c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", tt.id)})

			module.moduleName.controller.Update(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %d; want %d", w.Code, tt.expectedStatus)
			}

			if !tt.shouldFail {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName update successful!" {
					t.Error("expected success message in response, instead got:", response)
				}
			}
		})
	}
}

func TestModuleNameController_Delete(t *testing.T) {
	// Create a test moduleName
	created, err := module.moduleName.service.Create(&moduleName.CreateModuleName{
		Name: "Delete Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	tests := []struct {
		name           string
		id             uint
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:           "successful delete",
			id:             created.ID,
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "delete non-existent moduleName",
			id:             99999,
			expectedStatus: http.StatusOK, // GORM doesn't error on deleting non-existent records
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := testutil.CreateTestContext("DELETE", fmt.Sprintf("/moduleNames/%d", tt.id), nil)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", tt.id)})

			module.moduleName.controller.Delete(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %d; want %d", w.Code, tt.expectedStatus)
			}

			if !tt.shouldFail {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName deletion successful!" {
					t.Error("expected success message in response, instead got:", response)
				}
			}
		})
	}
}
