package gquery

import (
	"Gameserver-Query/games"
	"fmt"
)

type ServerQuery interface {
	QueryServer(address string, port int) (string, error)
}

var registry = map[string]ServerQuery{
	"minecraft": &games.MinecraftHandler{},
}

func Query(gameSlug string, address string, port int) (string, error) {
	handler, exists := registry[gameSlug]
	if !exists {
		return "", fmt.Errorf("Game: `%s` is not supported.", gameSlug)
	}

	queryResponse, err := handler.QueryServer(address, port)
	if err != nil {
		return "", err
	}

	return queryResponse, nil
}
