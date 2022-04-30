package mem

import (
	"sync"

	"github.com/ncomet/testcontainers-go/domain"
)

type AllGames struct {
	db sync.Map
}

func NewAllGames() *AllGames {
	return &AllGames{db: sync.Map{}}
}

func (a *AllGames) All() (allGames []*domain.Game) {
	a.db.Range(func(key, value interface{}) bool {
		allGames = append(allGames, value.(*domain.Game))
		return true
	})
	return
}

func (a *AllGames) Add(game *domain.Game) {
	a.db.Store(game.Id, game)
}

func (a *AllGames) Remove(game *domain.Game) {
	a.db.Delete(game.Id)
}

func (a *AllGames) By(id domain.GameId) *domain.Game {
	if game, present := a.db.Load(id); present {
		return game.(*domain.Game)
	}
	return nil
}
