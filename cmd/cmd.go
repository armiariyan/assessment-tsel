package cmd

import (
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"
	"github.com/armiariyan/assessment-tsel/internal/server"
)

func Run() {
	server.StartService(container.New())
}
