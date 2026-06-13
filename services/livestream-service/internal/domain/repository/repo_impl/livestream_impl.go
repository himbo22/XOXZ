package repo_impl

import (
	"context"
	"errors"
	"time"

	_const "github.com/himbo22/xoxz/livestream-service/internal/const"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/entity"
	"github.com/himbo22/xoxz/livestream-service/internal/domain/repository"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type livestreamRepository struct {
	collection *mongo.Collection
}

func NewLivestreamRepository(db *mongo.Database) repository.LivestreamRepository {
	return &livestreamRepository{
		collection: db.Collection(string(_const.LiveStreamRoom)),
	}
}

func (l *livestreamRepository) Create(ctx context.Context, s *entity.LivestreamRoom) error {
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now

	res, err := l.collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		s.ID = oid
	}

	return nil
}

func (l *livestreamRepository) GetByID(ctx context.Context, id string) (*entity.LivestreamRoom, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var room entity.LivestreamRoom

	err = l.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&room)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &room, nil
}

func (l *livestreamRepository) UpdateStatus(ctx context.Context, id string, status model.StreamStatus) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		},
	}

	res, err := l.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("invalid id")
	}

	return nil
}
