package handler

import (
	remindersProto "github.com/justIGreK/Reminders/pkg/go/reminders"
	"google.golang.org/grpc"
)

type Handler struct {
	server      grpc.ServiceRegistrar
	reminders RemindersService
}

func NewHandler(grpcServer grpc.ServiceRegistrar, rmsSRV RemindersService) *Handler {
	return &Handler{server: grpcServer, reminders: rmsSRV}
}
func (h *Handler) RegisterServices() {
	h.registerTxService(h.server, h.reminders)
}

func (h *Handler) registerTxService(server grpc.ServiceRegistrar, rms RemindersService) {
	remindersProto.RegisterRemindersServiceServer(server, &RemindersServiceServer{RmsSRV: rms})
}
