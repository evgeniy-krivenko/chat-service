package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	"github.com/evgeniy-krivenko/chat-service/internal/server"
	managerevents "github.com/evgeniy-krivenko/chat-service/internal/server-manager/events"
	managerv1 "github.com/evgeniy-krivenko/chat-service/internal/server-manager/v1"
	"github.com/evgeniy-krivenko/chat-service/internal/server/errhandler"
	eventstream "github.com/evgeniy-krivenko/chat-service/internal/services/event-stream"
	managerload "github.com/evgeniy-krivenko/chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/evgeniy-krivenko/chat-service/internal/services/manager-pool/in-mem"
	"github.com/evgeniy-krivenko/chat-service/internal/services/outbox"
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	canreceiveproblems "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/can-receive-problems"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
	freehands "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
	getmanagerprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-manager-profile"
	managerlogin "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/manager-login"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message"
	websocketstream "github.com/evgeniy-krivenko/chat-service/internal/websocket-stream"
)

const nameServerManager = "server-manager"

func initServerManager(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	keycloakClient *keycloakclient.Client,

	resource string,
	role string,
	secWSProtocol string,

	mLoadSrv *managerload.Service,
	mPool *inmemmanagerpool.Service,
	stream eventstream.EventStream,
	chatsRepo *chatsrepo.Repo,
	msgRepo *messagesrepo.Repo,
	problemsRepo *problemsrepo.Repo,
	profilesRepo *profilesrepo.Repo,
	outboxSvc *outbox.Service,
	db *store.Database,

	isProduction bool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	canReceiveProblemUseCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(mLoadSrv, mPool))
	if err != nil {
		return nil, fmt.Errorf("create can receive problem use case: %v", err)
	}

	freeHandsUseCase, err := freehands.New(freehands.NewOptions(mLoadSrv, mPool))
	if err != nil {
		return nil, fmt.Errorf("create free hands use case: %v", err)
	}

	getChatsUseCase, err := getchats.New(getchats.NewOptions(chatsRepo))
	if err != nil {
		return nil, fmt.Errorf("create get chats use case: %v", err)
	}

	getChatHistoryUseCase, err := getchathistory.New(getchathistory.NewOptions(msgRepo, problemsRepo))
	if err != nil {
		return nil, fmt.Errorf("create get chat history use case: %v", err)
	}

	sendMessageUseCase, err := sendmessage.New(sendmessage.NewOptions(msgRepo, problemsRepo, outboxSvc, db))
	if err != nil {
		return nil, fmt.Errorf("create send message usecase: %v", err)
	}

	closeChatUseCase, err := closechat.New(closechat.NewOptions(problemsRepo, msgRepo, outboxSvc, db))
	if err != nil {
		return nil, fmt.Errorf("create close chat use case: %v", err)
	}

	loginUseCase, err := managerlogin.New(managerlogin.NewOptions(
		keycloakClient,
		&keycloakclient.UserGetter{},
		profilesRepo,

		resource,
		role,
	))
	if err != nil {
		return nil, fmt.Errorf("create manager login use case: %v", err)
	}

	getManagerProfileUseCase, err := getmanagerprofile.New(getmanagerprofile.NewOptions(profilesRepo))
	if err != nil {
		return nil, fmt.Errorf("create get manager profile use case: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(
		canReceiveProblemUseCase,
		freeHandsUseCase,
		getChatsUseCase,
		getChatHistoryUseCase,
		sendMessageUseCase,
		closeChatUseCase,
		loginUseCase,
		getManagerProfileUseCase,
	))
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

	shutdown := make(chan struct{})

	wsManager, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		zap.L().Named("websocket-manager"),
		stream,
		managerevents.Adapter{},
		websocketstream.JSONEventWriter{},
		websocketstream.NewUpgrader(allowOrigins, secWSProtocol),
		shutdown,
	))
	if err != nil {
		return nil, fmt.Errorf("init websocket manager: %v", err)
	}

	return server.New(server.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		keycloakClient,
		resource,
		role,
		errHandleFunc,
		func(router *echo.Group) {
			managerv1.RegisterHandlers(router, v1Handlers)
		},
		func() {
			close(shutdown)
		},
		wsManager,
	))
}
