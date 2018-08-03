package server

import (
	"context"
	"os/exec"
	"time"

	"github.com/sdeoras/scheduler"
	"github.com/sdeoras/wg/wg/proto"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	group  scheduler.Scheduler
}

func New(ctx context.Context) (*Server, context.Context) {
	s := new(Server)
	s.group, s.ctx = scheduler.New(ctx, scheduler.NewTimeoutTrigger(0))
	ctx, s.cancel = context.WithCancel(context.Background())
	return s, ctx
}

func (s *Server) Run(ctx context.Context, request *proto.RunRequest) (*proto.RunResponse, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return &proto.RunResponse{}, CtxWithoutDeadline
	}

	ctx, cancel := context.WithTimeout(ctx, deadline.Sub(time.Now()))
	defer cancel()

	s.group.Go(s.ctx, func() error {
		commands := request.Commands
		cmd := exec.Command(commands[0], commands[1:]...)
		err := cmd.Run()
		return err
	})

	return &proto.RunResponse{}, nil
}

func (s *Server) Wait(ctx context.Context, request *proto.WaitRequest) (*proto.WaitResponse, error) {
	err := s.group.Wait()
	defer s.cancel() // shuts down server
	return &proto.WaitResponse{Mesg: "wg server is shutting down"}, err
}
