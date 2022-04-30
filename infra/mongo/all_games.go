package mongo

import (
	"context"
	"log"
	"time"

	"github.com/ncomet/testcontainers-go/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AllGames struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewAllGames(client *mongo.Client) *AllGames {
	collection := client.Database("testing").Collection("games")
	return &AllGames{client: client, coll: collection}
}

func (a AllGames) All() []*domain.Game {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter, err := a.coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var allGamesDocs []bson.M
	if err = filter.All(ctx, &allGamesDocs); err != nil {
		log.Fatal(err)
	}
	var allGames []*domain.Game
	for _, game := range allGamesDocs {
		allGames = append(allGames, toDomain(game))
	}
	return allGames
}

func (a AllGames) Add(game *domain.Game) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := a.coll.ReplaceOne(ctx, bson.D{{"_id", game.Id}}, bson.D{
		{"title", game.Title},
		{"PEGI", game.PEGI},
	}, options.Replace().SetUpsert(true))
	if err != nil {
		log.Println(err)
	}
}

func (a AllGames) Remove(game *domain.Game) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = a.coll.DeleteOne(ctx, bson.D{{"_id", game.Id}})
}

func (a AllGames) By(id domain.GameId) *domain.Game {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res bson.M
	err := a.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}
	return toDomain(res)
}

func toDomain(res bson.M) *domain.Game {
	return &domain.Game{
		Id:    domain.GameId(res["_id"].(string)),
		Title: res["title"].(string),
		PEGI:  domain.PEGI(res["PEGI"].(int32)),
	}
}
