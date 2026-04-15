package moduleName_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dmawardi/goTemplate/internal/auth"
	"github.com/dmawardi/goTemplate/internal/cache"
	"github.com/dmawardi/goTemplate/internal/common"
	"github.com/dmawardi/goTemplate/internal/config"
	"github.com/dmawardi/goTemplate/internal/moduleName"
	"github.com/dmawardi/goTemplate/internal/testutil"
	"github.com/dmawardi/goTemplate/internal/user"
)

var module = &testModule{}

type moduleNameModule struct {
	repository moduleName.ModuleNameRepository
	service    moduleName.ModuleNameService
	controller moduleName.ModuleNameController
}

type testModule struct {
	moduleName moduleNameModule
}

var app config.AppConfig

// Slice of state setters for each module
var stateFuncs = []testutil.StateFunc{
	user.SetState,
	auth.SetState,
	moduleName.SetState,
}

// Initial setup before running tests in package
func TestMain(m *testing.M) {
	fmt.Printf("Setting up test connection\n")
	// Set URL in app state
	app.BaseURL = common.BuildBaseUrl()

	// Setup DB
	dbClient := testutil.SetupTestDB(&moduleName.ModuleName{})
	// Set Gorm client
	app.DbClient = dbClient
	// Setup new cache
	app.Cache = &cache.CacheMap{}

	// Set app state
	testutil.SetAppState(&app, stateFuncs)

	// Create repository
	module.moduleName.repository = moduleName.NewModuleNameRepository(dbClient)
	// Create service
	module.moduleName.service = moduleName.NewModuleNameService(module.moduleName.repository)
	// Create controller
	module.moduleName.controller = moduleName.NewModuleNameController(module.moduleName.service)

	// Run the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}
