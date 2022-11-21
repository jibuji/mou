/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */

package mou

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReplaceOneObject struct {
	ctx         context.Context
	coll        *mongo.Collection
	opts        *options.ReplaceOptions
	filter      interface{}
	replacement interface{}
}

func (o *ReplaceOneObject) Where(filter interface{}) *ReplaceOneObject {
	o.filter = filter
	return o
}

func (o *ReplaceOneObject) Upsert() *ReplaceOneObject {
	o.opts.SetUpsert(true)
	return o
}

func (o *ReplaceOneObject) Replace(doc interface{}) *ReplaceOneObject {
	o.replacement = doc
	return o
}

func (o *ReplaceOneObject) Do() (*mongo.UpdateResult, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	return o.coll.ReplaceOne(o.ctx, o.filter, o.replacement, o.opts)
}
