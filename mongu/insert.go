/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */
package mongu

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneObject struct {
	ctx  context.Context
	coll *mongo.Collection
}

func (o *InsertOneObject) Do(document interface{}) (ID string, err error) {
	rst, err := o.coll.InsertOne(o.ctx, document)
	if err != nil {
		return "", err
	}
	if rst.InsertedID != nil {
		return rst.InsertedID.(primitive.ObjectID).Hex(), err
	}
	return "", nil
}

type InsertManyObject struct {
	ctx  context.Context
	coll *mongo.Collection
	opts *options.InsertManyOptions
}

//If true, no writes will be executed after one fails. The default value is true.
func (o *InsertManyObject) Ordered(ordered bool) *InsertManyObject {
	o.opts.SetOrdered(!ordered)
	return o
}

func (o *InsertManyObject) Do(documents []interface{}) (*mongo.InsertManyResult, error) {
	return o.coll.InsertMany(o.ctx, documents, o.opts)
}
