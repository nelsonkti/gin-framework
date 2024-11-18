package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/judwhite/go-svc"
	"go-framework/config"
	"go-framework/internal"
	"go-framework/internal/router"
	"go-framework/internal/server"
	"go-framework/util/binder"
	"go-framework/util/xconfig"
	"go-framework/util/xlog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

var confFile = flag.String("file", "", "input file path")

type logicProgram struct {
	once       sync.Once
	svcContext *server.SvcContext
}

func main() {
	p := &logicProgram{}
	if err := svc.Run(p, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL); err != nil {
		fmt.Println(err)
	}
}

// svc 服务运行框架 程序启动时执行Init+Start, 服务终止时执行Stop
func (p *logicProgram) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *logicProgram) Start() error {
	flag.Parse()

	var c config.Conf
	xconfig.New(&c, *confFile)

	logger := xlog.NewLogger(c.Log.Path, c.App.Name)

	p.svcContext = server.NewSvcContext(c, logger)

	appCxt := internal.Register(p.svcContext)

	go func() {
		newApp(c, appCxt)
	}()

	return nil
}

func newApp(c config.Conf, appCxt *internal.AppContent) {
	// 创建并配置验证器
	r := gin.New()

	binding.Validator = new(binder.Validator)

	r.Use(otelgin.Middleware(c.App.Name, otelgin.WithPropagators(propagation.TraceContext{})))

	router.Register(r, appCxt)

	err := r.Run(c.Server.Http.Addr)

	if err != nil {
		panic(err)
	}
}

func (p *logicProgram) Stop() error {
	p.once.Do(func() {
		//defer p.svcContext.RedisClient.Close()
		//defer p.svcContext.DBEngine.Close()
	})
	return nil
}
