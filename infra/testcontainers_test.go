package infra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/ncomet/testcontainers-go/domain"
	"github.com/ncomet/testcontainers-go/infra/mem"
	"github.com/ncomet/testcontainers-go/infra/mongo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	mdriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type image struct {
	name string
	port string
}

var (
	allGames domain.AllGames

	mongoDB = image{
		name: "mongo:5.0.8",
		port: "27017",
	}
)

type container struct {
	testcontainers.Container
	URI string
}

func TestMain(m *testing.M) {
	var code = 1
	defer func() { os.Exit(code) }()

	ctx := context.Background()
	mongoContainer, err := setup(ctx, mongoDB)
	if err != nil {
		log.Printf("Unexpected error, fallback to mem repository implementations.\nerror: %s.\n", err)
		allGames = mem.NewAllGames()
	} else {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		client, err := mdriver.Connect(ctx, options.Client().ApplyURI(mongoContainer.URI))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
		if err != nil {
			log.Printf("Could not connect to mongodb, fallback to mem repository implementations.\nerror: %s.\n", err)
			allGames = mem.NewAllGames()
		} else {
			allGames = mongo.NewAllGames(client)
		}
	}

	code = m.Run()
}

func setup(ctx context.Context, image image) (*container, error) {
	cont, uri, err := prepareContainer(ctx, image)
	if err != nil {
		return nil, err
	}
	return &container{Container: cont, URI: uri}, nil
}

func prepareContainer(ctx context.Context, image image) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        image.name,
		ExposedPorts: []string{image.port + "/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port(image.port)),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(image.port))
	if err != nil {
		return nil, "", err
	}

	var uri string
	switch image {
	case mongoDB:
		uri = fmt.Sprintf("mongodb://%s:%s", hostIP, mappedPort.Port())
	default:
		return nil, "", errors.New("TestContainers: unsupported image: " + image.name)
	}

	log.Printf("TestContainers: container %s is now running at %s\n", req.Image, uri)
	return container, uri, nil
}
