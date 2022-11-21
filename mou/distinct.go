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

type DistinctObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	opts   *options.DistinctOptions
	filter interface{}
	field  string
}

func (o *DistinctObject) Where(filter interface{}) *DistinctObject {
	o.filter = filter
	return o
}

func (o *DistinctObject) Field(fieldName string) *DistinctObject {
	o.field = fieldName
	return o
}

func (o *DistinctObject) Do() ([]interface{}, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	return o.coll.Distinct(o.ctx, o.field, o.filter)
}
