package postgres

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/thirteenths/message-processing-microservice/internal/domains"
)

type Postgres interface {
	CreateMessage(_ context.Context, message domains.Message) (int, error)
	UpdateStatusMessage(ctx context.Context, message domains.Message) error
	GetCountMessage(ctx context.Context) (count int, err error)
	GetProcessingCountMessage(ctx context.Context) (count int, err error)
}

type messageStorage struct {
	conn *pgx.Conn
}

func NewMessageStorage(url string) (Postgres, error) {
	conf, err := pgx.ParseConnectionString(url)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse connection string")
	}
	conn, err := pgx.Connect(conf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to connect to database")
	}
	return &messageStorage{conn: conn}, nil
}

const queryCreateMessage = "INSERT INTO MESSAGE(TEXT, KEY) VALUES ($1, $2) RETURNING ID"

func (s *messageStorage) CreateMessage(ctx context.Context, message domains.Message) (id int, err error) {
	err = s.conn.QueryRow(queryCreateMessage, message.Text, message.Key).Scan(&id)
	if err != nil {
		err = errors.WithMessage(err, "failed to insert message")
		return -1, err
	}
	return id, err
}

const queryUpdateStatusMessage = "UPDATE MESSAGE SET STATUS = TRUE WHERE KEY = $1 "

func (s *messageStorage) UpdateStatusMessage(ctx context.Context, message domains.Message) error {
	_, err := s.conn.Exec(queryUpdateStatusMessage, message.Key)
	if err != nil {
		err = errors.WithMessage(err, "failed to update status message")
	}

	return err
}

const queryGetCountMessage = "SELECT COUNT(STATUS) FROM MESSAGE"

func (s *messageStorage) GetCountMessage(ctx context.Context) (count int, err error) {
	err = s.conn.QueryRow(queryGetCountMessage).Scan(&count)
	if err != nil {
		err = errors.WithMessage(err, "failed to query count message")
		return 0, err
	}

	return count, err
}

const queryGetProcessingCountMessage = "SELECT COUNT(STATUS) FROM MESSAGE WHERE STATUS = FALSE"

func (s *messageStorage) GetProcessingCountMessage(ctx context.Context) (count int, err error) {
	err = s.conn.QueryRow(queryGetProcessingCountMessage).Scan(&count)
	if err != nil {
		err = errors.WithMessage(err, "failed to query count message")
		return 0, err
	}

	return count, err
}
