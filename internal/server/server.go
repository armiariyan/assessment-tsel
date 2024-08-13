package server

import (
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"
	"github.com/armiariyan/assessment-tsel/internal/server/http"
)

func StartService(container *container.Container) {
	http.StartH2CServer(container)
}
