package moduleName

import "github.com/dmawardi/goTemplate/internal/config"

var app *config.AppConfig

func SetState(appConfig *config.AppConfig) {
	app = appConfig
}
