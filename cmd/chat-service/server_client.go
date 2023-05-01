package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/server"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/server/errhandler"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	gethistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	keycloakClient *keycloakclient.Client,
	resource string,
	role string,

	msgRepo *messagesrepo.Repo,
	chatRepo *chatsrepo.Repo,
	problemRepo *problemsrepo.Repo,
	outboxSrv *outbox.Service,
	db *store.Database,

	isProduction bool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create get history usecase: %v", err)
	}

	sendMsgUseCase, err := sendmessage.New(sendmessage.NewOptions(chatRepo, msgRepo, outboxSrv, problemRepo, db))
	if err != nil {
		return nil, fmt.Errorf("create send message usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase, sendMsgUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(
		lg,
		isProduction,
		errhandler.ResponseBuilder,
	))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	errHandleFunc := errHandler.Handle

	srv, err := server.New(server.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		keycloakClient,
		resource,
		role,
		errHandleFunc,
		func(router server.EchoRouter) {
			clientv1.RegisterHandlers(router, v1Handlers)
		},
	))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", nameServerClient, err)
	}

	return srv, nil
}
