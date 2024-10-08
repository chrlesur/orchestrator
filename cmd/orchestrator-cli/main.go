package main

import (
    "fmt"
    "os"
    "orchestrator/cmd/orchestrator-cli/commands"
)

func main() {
    if err := commands.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}