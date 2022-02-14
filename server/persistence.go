package server

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type persister interface {
	get(ctx context.Context, key string) (*entity, error)
	set(ctx context.Context, key string, value *entity) error
	exists(ctx context.Context, key string) (bool, error)
}

type persistence struct {
	client     *firestore.Client
	collection string
}

func newPersistence(collection string) (*persistence, error) {
	opt := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(context.Background())

	return &persistence{
		client:     client,
		collection: collection,
	}, nil
}

func (p persistence) get(ctx context.Context, key string) (*entity, error) {
	docRef := p.client.Collection(p.collection).Doc(key)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	if status.Code(err) == codes.NotFound {
		return nil, nil
	}

	var e entity
	docSnap.DataTo(&e)
	return &e, nil

}

func (p persistence) set(ctx context.Context, key string, value *entity) error {
	doc := p.client.Collection(p.collection).Doc(key)
	_, err := doc.Set(ctx, value)

	if err != nil {
		return err
	}
	return nil
}

func (p persistence) exists(ctx context.Context, key string) (bool, error) {
	doc := p.client.Collection(p.collection).Doc(key)
	_, err := doc.Get(ctx)

	if status.Code(err) == codes.NotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
