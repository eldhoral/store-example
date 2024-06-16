package cmd

import (
    "os"
    "os/signal"
    "syscall"

    "store-api/app/api"

    "github.com/pkg/errors"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

var HttpCmd = &cobra.Command{
    Use:   "http serve",
    Short: "Run Http API",
    Long:  "Run Http API",
    RunE: func(cmd *cobra.Command, args []string) error {
        logrus.Infof("Starting the server at :%s", os.Getenv("HTTP_SERVER_PORT"))
        initHTTP()

        app := api.New(os.Getenv("APP_NAME"), baseHandler, storeHandler)

        echan := make(chan error)
        go func() {
            echan <- app.Run()
        }()
        term := make(chan os.Signal, 1)
        signal.Notify(term, os.Interrupt, syscall.SIGTERM)

        select {
        case <-term:
            logrus.Infoln("signal terminated detected")
            return nil
        case err := <-echan:
            return errors.Wrap(err, "service runtime error")
        }
    },
}
