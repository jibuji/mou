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
)

type DeleteOneObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	filter interface{}
}

func (o *DeleteOneObject) Where(filter interface{}) *DeleteOneObject {
	o.filter = filter
	return o
}

func (o *DeleteOneObject) Do() (deleteCount int, err error) {
	rst, er := o.coll.DeleteOne(o.ctx, o.filter)
	err = er
	if rst != nil {
		deleteCount = int(rst.DeletedCount)
	}
	return
}

type DeleteManyObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	filter interface{}
}

func (o *DeleteManyObject) Where(filter interface{}) *DeleteManyObject {
	o.filter = filter
	return o
}

func (o *DeleteManyObject) Do() (deleteCount int, err error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	rst, er := o.coll.DeleteMany(o.ctx, o.filter)
	err = er
	if rst != nil {
		deleteCount = int(rst.DeletedCount)
	}
	return
}
