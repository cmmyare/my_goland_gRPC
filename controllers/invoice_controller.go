package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/cmmyare24/go-gRPC/invoicer"
	"github.com/cmmyare24/go-gRPC/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateInvoice(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.CreateResponse, error) {
    invoice := models.Invoice{
        Amount: models.Amount{
            Amount:   req.Amount.Amount,
            Currence: req.Amount.Currence,
        },
        From:        req.From,
        To:          req.To,
        Description: req.Description,
    }

    _, err := models.Collection.InsertOne(ctx, invoice)
    if err != nil {
        log.Printf("Mongo insert error: %v", err)
        return nil, err
    }

    return &invoicer.CreateResponse{
        Amount:      req.Amount,
        To:          req.To,
        From:        req.From,
        Description: req.Description,
    }, nil
}


func UpdateInvoice(ctx context.Context, req *invoicer.UpdateRequest) (bool, error) {
    objID, err := primitive.ObjectIDFromHex(req.Id)
    if err != nil {
        log.Printf("Invalid ID: %v", err)
        return false, err
    }

    updateFields := bson.M{}

    if req.Amount != nil {
        updateFields["amount.amount"] = req.Amount.Amount
        updateFields["amount.currence"] = req.Amount.Currence
    }
    if req.From != "" {
        updateFields["from"] = req.From
    }
    if req.To != "" {
        updateFields["to"] = req.To
    }
    if req.Description != "" {
        updateFields["description"] = req.Description
    }

    if len(updateFields) == 0 {
        return false, fmt.Errorf("no fields provided to update")
    }

    filter := bson.M{"_id": objID}
    update := bson.M{"$set": updateFields}

    result, err := models.Collection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Printf("Update error: %v", err)
        return false, err
    }

    return result.ModifiedCount > 0, nil
}