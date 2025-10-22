package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Invoice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Amount      Amount             `bson:"amount"`
	From        string             `bson:"from"`
    To          string             `bson:"to"`
	Description string             `bson:"description"`
}

type Amount struct {
	Amount   int64  `bson:"amount"`
	Currence string `bson:"currence"`
}
