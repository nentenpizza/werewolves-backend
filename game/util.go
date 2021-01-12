package game

import (
	"errors"
	"fmt"
	"math/rand"
)

func rolesList(playerCount int) ([]func() Character, error) {
	if playerCount > MaxPlayers {
		return nil, errors.New("game: players in room must be <= 10")
	}
	roles, ok := rolesMap[playerCount]
	if !ok {
		return nil, fmt.Errorf("game: rolesMap for playerCount %d does not exists", playerCount)
	}
	var dst = make([]func() Character, len(roles), cap(roles))
	for k, v := range roles {
		dst[k] = v
	}
	rand.Shuffle(len(dst), func(i, j int) { dst[i], dst[j] = dst[j], dst[i] })
	return dst, nil
}
