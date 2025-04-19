package tgservice

import (
	"context"
	"strconv"

	"tg_service/telegrampb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"ext_service/internal/config"
)

type Client struct {
	cfg config.TGService
	cli telegrampb.MembershipServiceClient
}

func NewClient(ctx context.Context, cfg config.TGService) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Host+":"+strconv.Itoa(cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cli := telegrampb.NewMembershipServiceClient(conn)
	return &Client{
		cfg: cfg,
		cli: cli,
	}, nil
}

func (c *Client) ChekMembership(ctx context.Context, botToken, channelUrl string, userID int64) (*telegrampb.CheckResponse, error) {
	resp, err := c.cli.CheckMembership(ctx, &telegrampb.CheckRequest{
		BotToken:   botToken,
		ChannelUrl: channelUrl,
		UserId:     userID,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
