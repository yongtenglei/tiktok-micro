package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yunyandz.com/tiktok/registration"
	"yunyandz.com/tiktok/user-part/settings"
	"yunyandz.com/tiktok/user-part/web/controller"
	"yunyandz.com/tiktok/user-part/web/routers"
)

var (
	register registration.ConsulRegister
	id       string
)

func init() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)

	register = registration.NewConsulRegister(
		settings.ServiceConf.ConsulConf.Host,
		settings.ServiceConf.ConsulConf.Port)

	id = uuid.New().String()

	err := register.RegisterCheckByHTTP(
		settings.ServiceConf.UserWebClientConf.Name,
		id,
		settings.ServiceConf.UserWebClientConf.Host,
		settings.ServiceConf.UserWebClientConf.Port,
		settings.ServiceConf.UserWebClientConf.Tags,
	)

	if err != nil {
		panic(err)
	}

	controller.InitUserServiceClient()
}

func main() {

	router := routers.SetUp()

	addr := fmt.Sprintf("%s:%d",
		settings.ServiceConf.UserWebClientConf.Host,
		settings.ServiceConf.UserWebClientConf.Port)

	fmt.Println("User Web Client running on ", addr)
	if err := router.Run(addr); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig
	// Deregister
	//internal.ConsulDeRegister(setting.ProductServiceConf.ProductWebClientConfig.ID)
	err := register.DeRegister(id)
	if err != nil {
		zap.S().Errorw(fmt.Sprintf("Register DeRegister Consul failed on %s", addr), "err", err.Error())
		fmt.Printf("Deregister failed, %s\n", err.Error())
	} else {
		fmt.Println("Deregister ok")
	}

	zap.S().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server Shutdown", zap.Error(err))
	}

	zap.S().Info("Server exiting")

}
