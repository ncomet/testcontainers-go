package infra

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ncomet/testcontainers-go/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	gameId := domain.GameId(uuid.NewString())

	allGames.Add(&domain.Game{
		Id:    gameId,
		Title: "Assassin's Creed Valhalla",
		PEGI:  domain.Eighteen,
	})

	game := allGames.By(gameId)
	assert.Equal(t, gameId, game.Id)
	assert.Equal(t, "Assassin's Creed Valhalla", game.Title)
	assert.Equal(t, domain.Eighteen, game.PEGI)
}

func Test_Replace(t *testing.T) {
	gameId := domain.GameId(uuid.NewString())
	allGames.Add(&domain.Game{
		Id:    gameId,
		Title: "Rayman Raving Rabbids",
		PEGI:  domain.Seven,
	})

	allGames.Add(&domain.Game{
		Id:    gameId,
		Title: "Rainbow Six Siege",
		PEGI:  domain.Eighteen,
	})

	game := allGames.By(gameId)
	assert.Equal(t, gameId, game.Id)
	assert.Equal(t, "Rainbow Six Siege", game.Title)
	assert.Equal(t, domain.Eighteen, game.PEGI)
}

func Test_By(t *testing.T) {
	gameId := domain.GameId(uuid.NewString())
	allGames.Add(&domain.Game{
		Id:    gameId,
		Title: "Just Dance 2022",
		PEGI:  domain.Three,
	})
	allGames.Add(&domain.Game{
		Id:    domain.GameId(uuid.NewString()),
		Title: "Far Cry 6",
		PEGI:  domain.Eighteen,
	})

	game := allGames.By(gameId)

	assert.Equal(t, gameId, game.Id)
	assert.Equal(t, "Just Dance 2022", game.Title)
	assert.Equal(t, domain.Three, game.PEGI)
}

func Test_All(t *testing.T) {
	previous := len(allGames.All())
	allGames.Add(&domain.Game{
		Id:    domain.GameId(uuid.NewString()),
		Title: "Far Cry 6",
		PEGI:  domain.Eighteen,
	})

	assert.Equal(t, previous+1, len(allGames.All()))
}
