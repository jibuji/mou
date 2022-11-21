/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */

package mou

import (
	"context"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateOneObject struct {
	ctx     context.Context
	coll    *mongo.Collection
	opts    *options.UpdateOptions
	filter  interface{}
	updater interface{}
}

func (o *UpdateOneObject) Where(filter interface{}) *UpdateOneObject {
	o.filter = filter
	return o
}

func (o *UpdateOneObject) Upsert() *UpdateOneObject {
	o.opts.SetUpsert(true)
	return o
}

func (o *UpdateOneObject) Set(doc interface{}) *UpdateOneObject {
	var updater bson.M
	if o.updater == nil {
		updater = bson.M{}
	} else {
		updater = o.updater.(bson.M)
	}

	updater["$set"] = doc
	o.updater = updater
	return o
}

func (o *UpdateOneObject) Unset(doc interface{}) *UpdateOneObject {
	var updater bson.M
	if o.updater == nil {
		updater = bson.M{}
	} else {
		updater = o.updater.(bson.M)
	}

	updater["$unset"] = doc
	o.updater = updater
	return o
}

func (o *UpdateOneObject) Pull(doc interface{}) *UpdateOneObject {
	var updater bson.M
	if o.updater == nil {
		updater = bson.M{}
	} else {
		updater = o.updater.(bson.M)
	}
	updater["$pull"] = doc
	o.updater = updater
	return o
}

func (o *UpdateOneObject) Updater(doc interface{}) *UpdateOneObject {
	o.updater = doc
	return o
}

// warning: 如果没有发生insert（upsert时，已经存在匹配的记录），则不会返回ID
func (o *UpdateOneObject) Do() (string, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	rst, err := o.coll.UpdateOne(o.ctx, o.filter, o.updater, o.opts)
	// log.Printf("rst=%+v, err:%+v", rst, err)
	if err != nil {
		return "", errors.Wrapf(err,
			"error happened on call mongo api: filter:%v, updater:%v", o.filter, o.updater)
	}
	if id, ok := rst.UpsertedID.(primitive.ObjectID); ok {
		return id.Hex(), err
	}
	if id, ok := rst.UpsertedID.(string); ok {
		return id, err
	}
	return "", err
}

type UpdateManyObject struct {
	ctx     context.Context
	coll    *mongo.Collection
	opts    *options.UpdateOptions
	filter  interface{}
	updater interface{}
}

func (o *UpdateManyObject) Upsert() *UpdateManyObject {
	o.opts.SetUpsert(true)
	return o
}

func (o *UpdateManyObject) Where(filter interface{}) *UpdateManyObject {
	o.filter = filter
	return o
}

func (o *UpdateManyObject) Set(doc interface{}) *UpdateManyObject {
	o.updater = bson.M{"$set": doc}
	return o
}

func (o *UpdateManyObject) Updater(doc interface{}) *UpdateManyObject {
	o.updater = doc
	return o
}

func (o *UpdateManyObject) Do() (*mongo.UpdateResult, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	return o.coll.UpdateMany(o.ctx, o.filter, o.updater, o.opts)
}
