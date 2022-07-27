package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"go_chat_tutorial/internal/server"
	"go_chat_tutorial/pkg/logger"
	"os"
	"os/signal"
)

func handleErr(err error, msg string) {
	if err != nil {
		logger.Tag("main").Errorf(msg, err)
		os.Exit(1)
	}
}

func main() {

	logger.Init("ws-gateway")            // log khởi tạo service
	interrupt := make(chan os.Signal, 1) // tạo 1 kênh bắt sự kiện đến service chuẩn chỉnh
	signal.Notify(interrupt, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sgl := <-interrupt
		logger.WithField("signal", sgl.String()).Debug("interrupted")

		cancel()
	}()

	appConfig := &server.Config{} // lấy các thông tin config
	handleErr(env.Parse(appConfig), "failed to parse app config: %v")
	hub := server.NewHub(appConfig, logger.Tag("hub"))
	handleErr(hub.Run(ctx), "failed to start app: %v") // run server
}
