package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/armiariyan/assessment-tsel/internal/infrastructure/container"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
	"github.com/armiariyan/assessment-tsel/internal/server/handler"
)

func StartH2CServer(container *container.Container) {
	e := echo.New()

	SetupMiddleware(e, container)
	handler.SetupRouter(e, container)

	e.Server.Addr = fmt.Sprintf("%s:%s", container.Config.Apps.Address, container.Config.Apps.HttpPort)

	color.Println(color.Green(fmt.Sprintf("⇨ h2c server started on port: %s\n", container.Config.Apps.HttpPort)))
	log.Info(context.Background(), "h2c server started on port: "+container.Config.Apps.HttpPort)

	// * HTTP/2 Cleartext Server (HTTP2 over HTTP)
	// ! grace http can not used if user use windows. But its better http serve because the gracefully feature.
	// ! if user in windows download the libs, there will be error so i will use http serve from echo
	// gracehttp.Serve(&http.Server{Addr: e.Server.Addr, Handler: h2c.NewHandler(e, &http2.Server{MaxConcurrentStreams: 500, MaxReadFrameSize: 1048576})})

	s := http.Server{
		Addr:    e.Server.Addr,
		Handler: h2c.NewHandler(e, &http2.Server{MaxConcurrentStreams: 500, MaxReadFrameSize: 1048576}),
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(context.Background(), "failed echo listen and serve", err)
	}

}
