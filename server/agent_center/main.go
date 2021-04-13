package main

import (
	"fmt"
	"github.com/bytedance/Elkeid/server/agent_center/common"
	"github.com/bytedance/Elkeid/server/agent_center/grpctrans"
	"github.com/bytedance/Elkeid/server/agent_center/httptrans"
	"github.com/bytedance/Elkeid/server/agent_center/svr_registry"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
)

func init() {
	signal.Notify(common.Sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {
	fmt.Println("[MAIN] START_SERVER")

	//start http server and grpc server
	go httptrans.Run()
	go grpctrans.Run()

	//start pprof for debug
	go debug()

	//register to service discovery center
	regGrpc := svr_registry.NewGRPCServerRegistry()
	defer func() {
		regGrpc.Stop()
	}()

	regHttp := svr_registry.NewHttpServerRegistry()
	defer func() {
		regHttp.Stop()
	}()

	<-common.Sig
}

func debug() {
	//start pprof for debug
	if common.PProfEnable {
		err := http.ListenAndServe(fmt.Sprintf(":%d", common.PProfPort), nil)
		if err != nil {
			fmt.Printf("[MAIN] pprof ListenAndServe Error %s", err.Error())
		}
	}
}
