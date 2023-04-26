package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	serverclient "github.com/evgeniy-krivenko/chat-service/internal/server-client"
	"github.com/evgeniy-krivenko/chat-service/internal/server-client/errhandler"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
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
	db *store.Database,
	isProduction bool,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create get history usecase: %v", err)
	}
	sendMsgUseCase, err := sendmessage.New(sendmessage.NewOptions(chatRepo, msgRepo, problemRepo, db))
	if err != nil {
		return nil, fmt.Errorf("create send message usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg, getHistoryUseCase, sendMsgUseCase))
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

	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		v1Handlers,
		keycloakClient,
		resource,
		role,
		errHandleFunc,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
