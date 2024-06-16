package cmd

import (
	"os"
	"path"
	"runtime"
)

func initParams() map[string]string {
	params = make(map[string]string)

	params["APP_ENV"] = os.Getenv("APP_ENV")
	params["app-version"] = os.Getenv("APP_VERSION")
	params["app-name"] = os.Getenv("APP_NAME")

	_, b, _, _ := runtime.Caller(0)
	appDir := path.Join(path.Dir(b), "..")
	params["app-dir"] = appDir

	return params
}
