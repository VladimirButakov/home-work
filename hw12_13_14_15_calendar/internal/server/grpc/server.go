package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/VladimirButakov/home-work/tree/master/hw12_13_14_15_calendar/internal/app"
	gw "github.com/VladimirButakov/home-work/tree/master/hw12_13_14_15_calendar/internal/server/eventpb/EventService"
	"github.com/VladimirButakov/home-work/tree/master/hw12_13_14_15_calendar/internal/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	app    app.App
	server *http.Server
}

type grpcserver struct {
	gw.UnimplementedCalendarServer
	app app.App
}

var ErrBadRequest = errors.New("bad request")

func NewServer(app *app.App, address string, port string, grpcPort string) (*Server, error) {
	grpcServerEndpoint := net.JoinHostPort(address, grpcPort)

	lis, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to listen, %w", err)
	}

	logger := app.GetLogger().GetInstance()

	s := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_zap.StreamServerInterceptor(logger),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)))

	gw.RegisterCalendarServer(s, &grpcserver{app: *app})

	go func() {
		err := s.Serve(lis)
		if err != nil {
			app.GetLogger().Error(fmt.Errorf("cannot serve grpc, %w", err).Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		grpcServerEndpoint,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server, %w", err)
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	gwmux := runtime.NewServeMux()
	err = gw.RegisterCalendarHandler(ctx, gwmux, conn)
	if err != nil {
		return nil, fmt.Errorf("cannot register calendar handler, %w", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	server := &http.Server{
		Addr:         net.JoinHostPort(address, port),
		Handler:      gwmux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{*app, server}, nil
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("cannot start gateway server, %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot shutdown gateway server, %w", err)
	}

	return nil
}

func (s *grpcserver) CreateEvent(ctx context.Context, in *gw.ShortEvent) (*gw.Message, error) {
	if in.Title == "" || in.Date == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot create event, %s", ErrBadRequest)
	}

	err := s.app.CreateEvent(in.Title, in.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create event, %s", err)
	}

	return &gw.Message{Message: "Created new"}, nil
}

func (s *grpcserver) UpdateEvent(ctx context.Context, in *gw.Event) (*gw.Message, error) {
	event := storage.Event{
		ID:           in.Id,
		Title:        in.Title,
		EventDate:    in.EventDate,
		Duration:     in.Duration,
		Description:  in.Description,
		UserID:       in.UserID,
		NoticeBefore: in.NoticeBefore,
	}

	err := s.app.UpdateEvent(event)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update event, %s", err)
	}

	return &gw.Message{Message: "Updated"}, nil
}

func (s *grpcserver) DeleteEvent(ctx context.Context, in *gw.ID) (*gw.Message, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "cannot delete event, %s", ErrBadRequest)
	}

	err := s.app.RemoveEvent(in.Id)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "event not found, %s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete event, %s", err)
	}

	return &gw.Message{Message: "Deleted"}, nil
}

func (s *grpcserver) TimeEvents(ctx context.Context, in *gw.Date, events []storage.Event) (*gw.Events, error) {
	gwEvents := &gw.Events{}
	gwEvents.Results = make([]*gw.Event, len(events))

	for i, event := range events {
		gwEvents.Results[i] = &gw.Event{
			Id:           event.ID,
			Title:        event.Title,
			EventDate:    event.EventDate,
			Duration:     event.Duration,
			Description:  event.Description,
			UserID:       event.UserID,
			NoticeBefore: event.NoticeBefore,
		}
	}

	return gwEvents, nil
}

func (s *grpcserver) DailyEvents(ctx context.Context, in *gw.Date) (*gw.Events, error) {
	if in.Date == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get daily events, %s", ErrBadRequest)
	}

	var time time.Time = time.Unix(in.Date, 0)

	events, err := s.app.DailyEvents(time)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get events, %s", err)
	}

	return s.TimeEvents(ctx, in, events)
}

func (s *grpcserver) WeeklyEvents(ctx context.Context, in *gw.Date) (*gw.Events, error) {
	if in.Date == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get weekly events, %s", ErrBadRequest)
	}

	var time time.Time = time.Unix(in.Date, 0)

	events, err := s.app.WeeklyEvents(time)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get events, %s", err)
	}

	return s.TimeEvents(ctx, in, events)
}

func (s *grpcserver) MonthEvents(ctx context.Context, in *gw.Date) (*gw.Events, error) {
	if in.Date == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot get month events, %s", ErrBadRequest)
	}

	var time time.Time = time.Unix(in.Date, 0)

	events, err := s.app.MonthEvents(time)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get events, %s", err)
	}

	return s.TimeEvents(ctx, in, events)
}
