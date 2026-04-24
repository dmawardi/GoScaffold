package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/goTemplate/internal/moduleName"
	"github.com/dmawardi/goTemplate/internal/testutil"
)

func TestModuleNameRoutes_FindAll(t *testing.T) {
	// Set up the test environment
	created1, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 1",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}
	created2, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 2",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}
	created3, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Test ModuleName 3",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	// Store expected moduleName details in a map for easy lookup by ID
	expectedModuleNameDetails := map[uint]*moduleName.ModuleName{
		created1.ID: created1,
		created2.ID: created2,
		created3.ID: created3,
	}

	t.Cleanup(func() {
		testModule.moduleName.Repository.BulkDelete([]int{int(created1.ID), int(created2.ID), int(created3.ID)})
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{name: "Find all moduleNames successfully", token: testUsers["admin"].token, expectedStatus: http.StatusOK},
		{name: "Find all moduleNames with user token", token: testUsers["user"].token, expectedStatus: http.StatusOK},
		{name: "Unauthorized access with invalid token", token: "invalid", expectedStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/api/moduleNames?limit=10", nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Add authorization header if token provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP request
			testRouter.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// If successful, verify response structure
			if tt.expectedStatus == http.StatusOK {
				var responseData map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&responseData); err != nil {
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

				// Verify expected moduleNames are present
				for _, moduleNameBeingSearched := range responseModuleNames {
					if expectedModuleNameDetails[moduleNameBeingSearched.ID] != nil {
						testutil.AssertFieldsEqual(t, moduleNameBeingSearched, *expectedModuleNameDetails[moduleNameBeingSearched.ID], "FindAll")
					}
				}
			}
		})
	}
}

func TestModuleNameRoutes_Find(t *testing.T) {
	// Create a test moduleName
	testModuleName, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Find Route Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	t.Cleanup(func() {
		testModule.moduleName.Repository.BulkDelete([]int{int(testModuleName.ID)})
	})

	tests := []struct {
		name           string
		moduleNameID   uint
		token          string
		expectedStatus int
		shouldFail     bool
	}{
		{name: "Find moduleName successfully with admin token", moduleNameID: testModuleName.ID, token: testUsers["admin"].token, expectedStatus: http.StatusOK, shouldFail: false},
		{name: "Find moduleName successfully with user token", moduleNameID: testModuleName.ID, token: testUsers["user"].token, expectedStatus: http.StatusOK, shouldFail: false},
		{name: "Find non-existent moduleName", moduleNameID: 99999, token: testUsers["admin"].token, expectedStatus: http.StatusBadRequest, shouldFail: true},
		{name: "Unauthorized access", moduleNameID: testModuleName.ID, token: "invalid", expectedStatus: http.StatusUnauthorized, shouldFail: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/moduleNames/%d", tt.moduleNameID), nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Add authorization header
			req.Header.Set("Authorization", "Bearer "+tt.token)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP request
			testRouter.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// If successful, verify response
			if !tt.shouldFail {
				var responseModuleName moduleName.ModuleName
				if err := json.NewDecoder(rr.Body).Decode(&responseModuleName); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				testutil.AssertFieldsEqual(t, responseModuleName, *testModuleName, tt.name)
			}
		})
	}
}

func TestModuleNameRoutes_Create(t *testing.T) {
	tests := []struct {
		name           string
		createData     *moduleName.CreateModuleName
		token          string
		expectedStatus int
		shouldFail     bool
	}{
		{
			name: "Create moduleName successfully with admin token",
			createData: &moduleName.CreateModuleName{
				Name: "Create Test ModuleName Admin",
			},
			token:          testUsers["admin"].token,
			expectedStatus: http.StatusCreated,
			shouldFail:     false,
		},
		{
			name: "Create moduleName with user token - forbidden",
			createData: &moduleName.CreateModuleName{
				Name: "Create Test ModuleName User",
			},
			token:          testUsers["user"].token,
			expectedStatus: http.StatusForbidden,
			shouldFail:     true,
		},
		{
			name: "Create moduleName with empty name",
			createData: &moduleName.CreateModuleName{
				Name: "",
			},
			token:          testUsers["admin"].token,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
		{
			name: "Create moduleName without token",
			createData: &moduleName.CreateModuleName{
				Name: "No Auth ModuleName",
			},
			token:          "invalid",
			expectedStatus: http.StatusUnauthorized,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request body
			bodyBytes, err := json.Marshal(tt.createData)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/moduleNames", bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Set content type and authorization header
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP request
			testRouter.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// If successful, verify response and clean up
			if !tt.shouldFail {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName creation successful!" {
					t.Error("expected success message in response, instead got:", response)
				}

				// Clean up created moduleName
				foundModuleName, err := testModule.moduleName.Service.FindByName(tt.createData.Name)
				if err == nil {
					testModule.moduleName.Repository.BulkDelete([]int{int(foundModuleName.ID)})
				}
			}
		})
	}
}

func TestModuleNameRoutes_Update(t *testing.T) {
	// Create a test moduleName to update
	testModuleName, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Update Route Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	t.Cleanup(func() {
		testModule.moduleName.Repository.BulkDelete([]int{int(testModuleName.ID)})
	})

	tests := []struct {
		name           string
		moduleNameID   uint
		updateData     *moduleName.UpdateModuleName
		token          string
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:         "Update moduleName successfully with admin token",
			moduleNameID: testModuleName.ID,
			updateData: &moduleName.UpdateModuleName{
				Name: "Updated Name Admin",
			},
			token:          testUsers["admin"].token,
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:         "Update moduleName with user token - forbidden",
			moduleNameID: testModuleName.ID,
			updateData: &moduleName.UpdateModuleName{
				Name: "Updated Name User",
			},
			token:          testUsers["user"].token,
			expectedStatus: http.StatusForbidden,
			shouldFail:     true,
		},
		{
			name:         "Update non-existent moduleName",
			moduleNameID: 99999,
			updateData: &moduleName.UpdateModuleName{
				Name: "Updated",
			},
			token:          testUsers["admin"].token,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
		{
			name:         "Update without token",
			moduleNameID: testModuleName.ID,
			updateData: &moduleName.UpdateModuleName{
				Name: "Unauthorized",
			},
			token:          "invalid",
			expectedStatus: http.StatusUnauthorized,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request body
			bodyBytes, err := json.Marshal(tt.updateData)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("PUT", fmt.Sprintf("/api/moduleNames/%d", tt.moduleNameID), bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Set content type and authorization header
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP request
			testRouter.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// If successful, verify response
			if !tt.shouldFail {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName update successful!" {
					t.Error("expected success message in response, instead got:", response)
				}
			}
		})
	}
}

func TestModuleNameRoutes_Delete(t *testing.T) {
	// Create a test moduleName to delete
	testModuleName, err := testModule.moduleName.Service.Create(&moduleName.CreateModuleName{
		Name: "Delete Route Test ModuleName",
	})
	if err != nil {
		t.Fatalf("Failed to create moduleName: %v", err)
	}

	tests := []struct {
		name           string
		moduleNameID   uint
		token          string
		expectedStatus int
	}{
		{name: "Delete moduleName successfully with admin token", moduleNameID: testModuleName.ID, token: testUsers["admin"].token, expectedStatus: http.StatusOK},
		{name: "Delete moduleName with user token - forbidden", moduleNameID: testModuleName.ID, token: testUsers["user"].token, expectedStatus: http.StatusForbidden},
		{name: "Delete non-existent moduleName", moduleNameID: 99999, token: testUsers["admin"].token, expectedStatus: http.StatusOK},
		{name: "Delete without token", moduleNameID: testModuleName.ID, token: "invalid", expectedStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/moduleNames/%d", tt.moduleNameID), nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Add authorization header
			req.Header.Set("Authorization", "Bearer "+tt.token)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve HTTP request
			testRouter.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// If successful, verify response
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if message, ok := response["message"]; !ok || message != "ModuleName deletion successful!" {
					t.Error("expected success message in response, instead got:", response)
				}
			}
		})
	}
}
