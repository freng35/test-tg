package service

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"tg_service/internal/config"
	"tg_service/store"
	"tg_service/telegrampb"
)

type Service struct {
	cfg          config.GRPC
	grpcServer   *grpc.Server
	pgRepository *store.PostgresRepository
}

func NewService(cfg config.GRPC, pgRepository *store.PostgresRepository) *Service {
	return &Service{
		cfg:          cfg,
		grpcServer:   nil,
		pgRepository: pgRepository,
	}
}

func (s *Service) CheckMembership(ctx context.Context, req *telegrampb.CheckRequest) (*telegrampb.CheckResponse, error) {
	bot, err := tgbotapi.NewBotAPI(req.BotToken)
	if err != nil {
		return nil, err
	}

	chatID, err := getChatIDFromURL(bot, req.ChannelUrl)
	if err != nil {
		return nil, err
	}

	isMember, err := checkTelegramMembership(bot, chatID, req.UserId)
	if err != nil {
		return nil, err
	}

	var added bool
	if isMember {
		err = s.pgRepository.AddUser(ctx, req.UserId, chatID)
		if err != nil {
			slog.Error("can not add user", slog.String("err", err.Error()))
		}
		added = err == nil
	}

	return &telegrampb.CheckResponse{
		IsMember:  isMember,
		AddedToDb: added,
	}, nil
}

func (s *Service) Run() error {
	s.makeServer()
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Close() {
	if s.grpcServer == nil {
		slog.Warn("grpc server not inited")
		return
	}

	s.grpcServer.GracefulStop()
	slog.Info("grpc server stopped")
}

func (s *Service) makeServer() {
	s.grpcServer = grpc.NewServer()
	reflection.Register(s.grpcServer)

	telegrampb.RegisterMembershipServiceServer(s.grpcServer, s)
}

func checkTelegramMembership(bot *tgbotapi.BotAPI, chatID, userID int64) (bool, error) {
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return false, err
	}

	return member.Status == "member" || member.Status == "administrator" || member.Status == "creator", nil
}

func getChatIDFromURL(bot *tgbotapi.BotAPI, channelURL string) (int64, error) {
	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			SuperGroupUsername: channelURL,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get chat ID: %w", err)
	}

	return chat.ID, nil
}
