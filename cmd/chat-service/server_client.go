package main

import (
	"fmt"
	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	getuserprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/server"
	clientevents "github.com/evgeniy-krivenko/chat-service/internal/server-client/events"
	clientv1 "github.com/evgeniy-krivenko/chat-service/internal/server-client/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/server/errhandler"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	gethistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/send-message"
	websocketstream "github.com/evgeniy-krivenko/chat-service/internal/websocket-stream"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	keycloakClient *keycloakclient.Client,

	resource string,
	role string,
	secWSProtocol string,

	msgRepo *messagesrepo.Repo,
	chatRepo *chatsrepo.Repo,
	problemRepo *problemsrepo.Repo,
	profilesRepo *profilesrepo.Repo,
	outboxSrv *outbox.Service,
	db *store.Database,
	stream eventstream.EventStream,

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
	getUserProfileUseCase, err := getuserprofile.New(getuserprofile.NewOptions(profilesRepo))
	if err != nil {
		return nil, fmt.Errorf("get user profile: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase, sendMsgUseCase, getUserProfileUseCase))
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

	shutdown := make(chan struct{})

	wsClient, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		zap.L().Named("websocket-client"),
		stream,
		clientevents.Adapter{},
		websocketstream.JSONEventWriter{},
		websocketstream.NewUpgrader(allowOrigins, secWSProtocol),
		shutdown,
	))
	if err != nil {
		return nil, fmt.Errorf("init websocket client: %v", err)
	}

	srv, err := server.New(server.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		keycloakClient,
		resource,
		role,
		errHandleFunc,
		func(router *echo.Group) {
			clientv1.RegisterHandlers(router, v1Handlers)
		},
		func() {
			close(shutdown)
		},
		wsClient,
	))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", nameServerClient, err)
	}

	return srv, nil
}
