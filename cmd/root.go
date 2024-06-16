package cmd

import (
    "os"

    "store-api/pkg/errs"

    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
    Use:   "Go Simple API",
    Short: "Go Simple API / Service Demo",
    Long:  "Go Simple API / Service Demo with HTTP API",
    Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

//register command
func init() {
    //load environment variable
    if err := godotenv.Load("./params/.env"); err != nil && err.Error() != errs.ErrNoSuchFile {
        logrus.Fatalln("unable to load environment variable", err.Error())
    }

    rootCmd.AddCommand(HttpCmd)
}

func Execute() error {
    cmd, _, err := rootCmd.Find(os.Args[1:])

    // Default run main http if not set
    if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
        args := append([]string{"http"}, os.Args[1:]...)
        rootCmd.SetArgs(args)
    }

    return rootCmd.Execute()
}
