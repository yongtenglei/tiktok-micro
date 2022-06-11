package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"yunyandz.com/tiktok/proto/pb"
	"yunyandz.com/tiktok/user-part/service/biz"
	"yunyandz.com/tiktok/user-part/settings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"yunyandz.com/tiktok/dao/mysql"
	"yunyandz.com/tiktok/logger"
	"yunyandz.com/tiktok/registration"
)

func main() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)
	settings.ParseConfig("setting.yaml")
	logger.New()
	mysql.Init()

	addr := fmt.Sprintf("%s:%d", settings.ServiceConf.UserWebServerConf.Host, settings.ServiceConf.UserWebServerConf.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Errorw("Create listener failed", "err", err.Error())
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &biz.UserServiceServer{})

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	register := registration.NewConsulRegister(
		settings.ServiceConf.ConsulConf.Host,
		settings.ServiceConf.ConsulConf.Port)

	id := uuid.New().String()
	err = register.RegisterCheckByGRPC(
		settings.ServiceConf.UserWebServerConf.Name,
		id,
		settings.ServiceConf.UserWebServerConf.Host,
		settings.ServiceConf.UserWebServerConf.Port,
		settings.ServiceConf.UserWebServerConf.Tags)
	if err != nil {
		zap.S().Errorw("Web Server register to Consul failed", "err", err.Error())
		panic(err)
	}

	fmt.Printf("User Server running on %s", addr)

	go func() {
		if err := server.Serve(listener); err != nil {
			zap.S().Errorw(fmt.Sprintf("Server Serve failed in %s", addr), "err", err.Error())
			panic(err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	// Deregister
	err = register.DeRegister(id)
	if err != nil {
		zap.S().Errorw(fmt.Sprintf("Register DeRegister Consul failed on %s", addr), "err", err.Error())
		fmt.Printf("Deregister failed, %s\n", err.Error())
	} else {
		fmt.Println("Deregister successfully")
	}

	zap.S().Info("Shutdown ...")

	server.GracefulStop()

	zap.S().Info("Server exiting")
}
