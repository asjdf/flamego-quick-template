package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/flamego/cors"
	"github.com/flamego/flamego"
	"github.com/spf13/cobra"
	"github.com/wujunyi792/gin-template-new/config"
	example "github.com/wujunyi792/gin-template-new/internal/app/example/router"
	file "github.com/wujunyi792/gin-template-new/internal/app/file/router"
	ws "github.com/wujunyi792/gin-template-new/internal/app/websocket/router"
	"github.com/wujunyi792/gin-template-new/internal/cache"
	"github.com/wujunyi792/gin-template-new/internal/database"
	"github.com/wujunyi792/gin-template-new/internal/middleware"
	"github.com/wujunyi792/gin-template-new/pkg/colorful"
	"github.com/wujunyi792/gin-template-new/pkg/ip"
	"github.com/wujunyi792/gin-template-new/pkg/logx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configYml string
	E         *flamego.Flame
	StartCmd  = &cobra.Command{
		Use:     "server",
		Short:   "Set Application config info",
		Example: "main server -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			println("Loading config...")
			setUp()
			println("Loading config complete")
			println("Loading Api...")
			load()
			println("Loading Api complete")
		},
		Run: func(cmd *cobra.Command, args []string) {
			println("Starting Server...")
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}

func setUp() {
	// 顺序不能变 logger依赖config logger后面的同时依赖logger和config 否则crash
	config.LoadConfig(configYml)
	logx.InitLogger()
	database.InitDB()
	cache.InitCache()
}

func load() {
	flamego.SetEnv(flamego.EnvType(config.GetConfig().MODE))
	E = flamego.New()
	E.Use(flamego.Recovery(), middleware.RequestLog(), flamego.Renderer(), cors.CORS())

	// TODO 新路由请在这注册
	example.AppExampleInit(E)
	file.AppFileInit(E)
	ws.AppWebsocketInit(E)
}

func run() {
	srv := &http.Server{
		Addr:    config.GetConfig().Host + ":" + config.GetConfig().Port,
		Handler: E,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			println(colorful.Red("Got Server Err: " + err.Error()))
		}
	}()

	println(colorful.Green("Server run at:"))
	println(fmt.Sprintf("-  Local:   http://localhost:%s", config.GetConfig().Port))
	for _, host := range ip.GetLocalHost() {
		println(fmt.Sprintf("-  Network: http://%s:%s", host, config.GetConfig().Port))
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	println(colorful.Blue("Shutting down server..."))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
	}

	println(colorful.Green("Server exiting Correctly"))
}