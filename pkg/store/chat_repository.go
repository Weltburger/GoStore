package store

import (
	"GoStore/pkg/models"
	"context"
	"encoding/json"
)

type ChatRepository struct {
	store *Store
}

func (chatRepository *ChatRepository) InsertChat(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	f, err := json.MarshalIndent(chat.Data, "", " ")
	if err != nil {
		return nil, err
	}

	if err := chatRepository.store.DB.QueryRowxContext(ctx, `insert into "public"."chats"(uuid, data) 
VALUES($1, $2) returning uuid`, chat.UUID, f).Scan(&chat.UUID); err != nil {
		return nil, err
	}

	return chat, nil
}
