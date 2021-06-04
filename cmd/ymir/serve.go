package ymir

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/svartlfheim/ymir/internal/server"
)

func serve(cmd YmirCommand) error {
	cfg := cmd.GetConfig()
	l := cmd.GetLogger()

	moduleRepo, err := buildModuleRepository(cmd.GetConfig(), cmd.cobra.Context(), l)

	if err != nil {
		l.Fatal().Err(err).Msg("failed to build module repository")
	}

	a, err := buildAuditor(cfg, cmd.cobra.Context(), l)

	if err != nil {
		l.Fatal().Err(err).Msg("failed to build auditor")
	}

	cb := buildCommandBus(cmd)
	h := server.NewServer([]server.Controller{
		&server.MiscController{},
		server.NewModulesController(l, cb, a),
		server.NewModuleRegistryController(l, moduleRepo, cb),
	})

	fmt.Printf("Listening on %s\n", cfg.Server.Port)
	err = http.ListenAndServe(":"+cfg.Server.Port, handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(os.Stdout, h)))

	return err

}
