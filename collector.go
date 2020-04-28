package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/n9e/win-collector/cache"
	"github.com/n9e/win-collector/config"
	"github.com/n9e/win-collector/sys/identity"

	"github.com/n9e/win-collector/http/routes"
	"github.com/n9e/win-collector/stra"
	"github.com/n9e/win-collector/sys"
	"github.com/n9e/win-collector/sys/funcs"
	"github.com/n9e/win-collector/sys/plugins"
	"github.com/n9e/win-collector/sys/ports"
	"github.com/n9e/win-collector/sys/procs"

	"github.com/didi/nightingale/src/toolkits/http"
	tlogger "github.com/didi/nightingale/src/toolkits/logger"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"
)

var (
	vers *bool
	help *bool
	conf *string
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	conf = flag.String("f", "", "specify configuration file.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", config.Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	aconf()
	pconf()
	start()
	cfg := config.Get()

	tlogger.Init(cfg.Logger)

	identity.Init(cfg.Identity)
	log.Println("endpoint:", identity.Identity)

	sys.Init(cfg.Sys)
	stra.Init(cfg.Stra)

	funcs.InitRpcClients()
	funcs.BuildMappers()
	funcs.Collect()

	//插件采集
	plugins.Detect()

	//进程采集
	procs.Detect()

	//端口采集
	ports.Detect()

	//初始化缓存，用作保存COUNTER类型数据
	cache.Init()

	r := gin.New()
	routes.Config(r)
	http.Start(r, "collector", cfg.Logger.Level)
	ending()
}

// auto detect configuration file
func aconf() {
	if *conf != "" && file.IsExist(*conf) {
		return
	}

	*conf = "etc/win-collector.local.yml"
	if file.IsExist(*conf) {
		return
	}

	*conf = "etc/win-collector.yml"
	if file.IsExist(*conf) {
		return
	}

	fmt.Println("no configuration file for collector")
	os.Exit(1)
}

// parse configuration file
func pconf() {
	if err := config.Parse(*conf); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
}

func ending() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())
	}

	logger.Close()
	http.Shutdown()
	fmt.Println("sender stopped successfully")
}

func start() {
	runner.Init()
	fmt.Println("collector start, use configuration file:", *conf)
	fmt.Println("runner.cwd:", runner.Cwd)
}
