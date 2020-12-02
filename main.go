package main

import (
	"github.com/epiphany-platform/m-azure-basic-infrastructure/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cmd.Execute()
}
