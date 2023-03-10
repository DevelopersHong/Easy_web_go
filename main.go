package main

import (
	web "Easy_web_go/pkg"
	static "Easy_web_go/web"
	_ "Easy_web_go/web/filters"
	"context"
	"fmt"
	"time"
)

func main() {
	shutdown := web.NewGracefulShutdown()
	server := web.NewSdkHttpServer("my-test-server",
		web.MetricFilterBuilder, shutdown.ShutdownFilterBuilder)
	adminServer := web.NewSdkHttpServer("admin-test-server",
		// 注意，如果你真实环境里面，使用的是多个 server监听不同端口，
		// 那么这个 shutdown最好也是多个。互相之间就不会有竞争
		// MetricFilterBuilder 是无状态的，所以不存在这种问题
		web.MetricFilterBuilder, shutdown.ShutdownFilterBuilder)

	// 注册路由
	_ = server.Route("POST", "/user/create/*", static.SignUp)
	_ = server.Route("POST", "/slowService", static.SlowService)

	// 准备静态路由

	staticHandler := web.NewStaticResourceHandler(
		"web/static", "/static",
		web.WithMoreExtension(map[string]string{
			"mp3": "audio/mp3",
		}), web.WithFileCache(1<<20, 100))
	// 访问 Get http://localhost:8080/static/forest.png
	server.Route("GET", "/static/*", staticHandler.ServeStaticResource)

	go func() {
		if err := adminServer.Start(":8081"); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := server.Start(":8080"); err != nil {
			// 快速失败，因为服务器都没启动成功，啥也做不了
			panic(err)
		}
		// 假设我们后面还有很多动作
	}()

	// 先执行 RejectNewRequestAndWaiting，等待所有的请求
	// 然后我们关闭 server，如果是多个 server，可以多个 goroutine 一起关闭
	//
	web.WaitForShutdown(
		func(ctx context.Context) error {
			// 假设我们这里有一个 hook
			// 可以通知网关我们要下线了
			fmt.Println("mock notify gateway")
			time.Sleep(time.Second * 2)
			return nil
		},
		shutdown.RejectNewRequestAndWaiting,
		// 全部请求处理完了我们就可以关闭 server了
		web.BuildCloseServerHook(server, adminServer),
		func(ctx context.Context) error {
			// 假设这里我要清理一些执行过程中生成的临时资源
			fmt.Println("mock release resources")
			time.Sleep(time.Second * 2)
			return nil
		})

	// filterNames := ReadFromConfig
	// 匿名引入之后，就可以在这里按名索引 filter
	//web.NewSdkHttpServerWithFilterNames("my-server", filterNames...)

}
