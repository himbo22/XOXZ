package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type LivestreamRoom struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomName    string        `bson:"room_name" json:"room_name"`
	PublisherID uuid.UUID     `bson:"publisher_id" json:"publisher_id"`
	Status      string        `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (LivestreamRoom) TableName() string {
	return "livestream_rooms"
}
