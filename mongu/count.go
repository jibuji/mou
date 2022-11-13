/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */

package mongu

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CountObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	opts   *options.CountOptions
	filter interface{}
}

func (o *CountObject) Where(filter interface{}) *CountObject {
	o.filter = filter
	return o
}

func (o *CountObject) Limit(i int64) *CountObject {
	o.opts.SetLimit(i)
	return o
}

func (o *CountObject) Skip(i int64) *CountObject {
	o.opts.SetSkip(i)
	return o
}

func (o *CountObject) Do() (int64, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	return o.coll.CountDocuments(o.ctx, o.filter, o.opts)
}
