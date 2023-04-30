package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/server"
	managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1"
	// managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1".
	"github.com/evgeniy-krivenko/chat-service/internal/server/errhandler"
	managerload "github.com/evgeniy-krivenko/chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/evgeniy-krivenko/chat-service/internal/services/manager-pool/in-mem"
	canreceiveproblems "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/free-hands"
)

const nameServerManager = "server-manager"

func initServerManager(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	keycloakClient *keycloakclient.Client,
	mLoadSrv *managerload.Service,
	resource string,
	role string,
	isProduction bool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	mPool := inmemmanagerpool.New()

	canReceiveProblemUseCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(mLoadSrv, mPool))
	if err != nil {
		return nil, fmt.Errorf("create can receive problem use case: %v", err)
	}

	freeHandsUseCase, err := freehands.New(freehands.NewOptions(mLoadSrv, mPool))
	if err != nil {
		return nil, fmt.Errorf("create free hands use case: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(canReceiveProblemUseCase, freeHandsUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 manager handlers: %v", err)
	}
	// other components
	errHandler, err := errhandler.New(errhandler.NewOptions(
		lg,
		isProduction,
		errhandler.ResponseBuilder,
	))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	errHandleFunc := errHandler.Handle

	return server.New(server.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		keycloakClient,
		resource,
		role,
		errHandleFunc,
		func(router server.EchoRouter) {
			managerv1.RegisterHandlers(router, v1Handlers)
		},
	))
}
