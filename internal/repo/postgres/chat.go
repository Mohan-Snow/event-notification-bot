package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"event-notification-bot/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

const (
	chatTable = "chats"
)

var (
	ErrChatNotFound   = errors.New("chat not found")
	ErrChatNotCreated = errors.New("chat not created")
)

func (p *Postgres) SaveChat(ctx context.Context, id int64) error {
	query := sq.Insert(chatTable).
		Columns("chat_id").
		Values(id).
		PlaceholderFormat(sq.Dollar)

	rowQuery, args, err := query.ToSql()
	if err != nil {
		log.Printf("repo.postgres.SaveChat1: %v", err)
		return err
	}

	commandTag, err := p.db.Exec(ctx, rowQuery, args...)
	if err != nil {
		fmt.Println(id)
		log.Printf("repo.postgres.SaveChat2: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrChatNotCreated
	}

	return err
}

func (p *Postgres) ListChats(ctx context.Context) ([]*model.Chat, error) {
	query := sq.Select("chat_id").
		From(chatTable)

	rowQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var orders []*model.Chat
	err = pgxscan.Select(ctx, p.db, &orders, rowQuery, args...)

	if err != nil {
		return nil, err
	}

	return orders, nil

}

func (p *Postgres) FindChatById(ctx context.Context, id int64) (*model.Chat, error) {
	query := sq.Select("chat_id").
		From(chatTable).
		Where(sq.Eq{"chat_id": id}).
		PlaceholderFormat(sq.Dollar)

	rowQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var order model.Chat
	err = pgxscan.Get(ctx, p.db, &order, rowQuery, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", ErrChatNotFound, err)
		}

		return nil, err
	}

	return &order, nil
}

func (p *Postgres) DeleteChat(ctx context.Context, id int64) error {
	query := sq.Delete(chatTable).
		Where(sq.Eq{"chat_id": id}).
		PlaceholderFormat(sq.Dollar)

	rowQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = p.db.Exec(ctx, rowQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
