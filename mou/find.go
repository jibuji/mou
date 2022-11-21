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

type FindManyObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	opts   *options.FindOptions
	filter interface{}
}

func (o *FindManyObject) Where(filter interface{}) *FindManyObject {
	o.filter = filter
	return o
}

func (o *FindManyObject) Skip(i int64) *FindManyObject {
	o.opts.SetSkip(i)
	return o
}

func (o *FindManyObject) Sort(sort interface{}) *FindManyObject {
	o.opts.SetSort(sort)
	return o
}

func (o *FindManyObject) Limit(i int64) *FindManyObject {
	o.opts.SetLimit(i)
	return o
}

func (o *FindManyObject) Select(projection interface{}) *FindManyObject {
	o.opts.Projection = projection
	return o
}

func (o *FindManyObject) Collation(c *options.Collation) *FindManyObject {
	o.opts.Collation = c
	return o
}

func (o *FindManyObject) Do(results interface{}) error {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	cursor, err := o.coll.Find(o.ctx, o.filter, o.opts)
	if err != nil {
		return err
	}
	defer cursor.Close(o.ctx)
	return cursor.All(o.ctx, results)
}

// TODO: 1. customize element type 2.
// 将来 golang 支持泛型后，可以用更优雅的方式改进这个函数
// 当错误发生或者context done，或者没有更多元素时，out chan会被自行关闭
func (o *FindManyObject) DoStream() (chan interface{}, error) {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	cursor, err := o.coll.Find(o.ctx, o.filter, o.opts)
	if err != nil {
		return nil, err
	}
	const BufferSize = 100
	outChan := make(chan interface{}, BufferSize)
	go func() {
	loop:
		for cursor.Next(o.ctx) {
			item := map[string]interface{}{}
			if err = cursor.Decode(&item); err != nil {
				// log.Errorf("decode item err:%v", err)
				break
			}
			// log.Printf("incoming value %#v", item)
			select {
			case <-o.ctx.Done():
				break loop
			case outChan <- item:
			}
		}
		close(outChan)
		cursor.Close(o.ctx)
	}()

	return outChan, nil
}

type FindOneObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	opts   *options.FindOneOptions
	filter interface{}
}

func (o *FindOneObject) Where(filter interface{}) *FindOneObject {
	o.filter = filter
	return o
}

func (o *FindOneObject) Skip(i int64) *FindOneObject {
	o.opts.SetSkip(i)
	return o
}

func (o *FindOneObject) Sort(sort interface{}) *FindOneObject {
	o.opts.SetSort(sort)
	return o
}

func (o *FindOneObject) Select(projection interface{}) *FindOneObject {
	o.opts.Projection = projection
	return o
}

func (o *FindOneObject) Do(result interface{}) error {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	rst := o.coll.FindOne(o.ctx, o.filter, o.opts)
	return rst.Decode(result)
}

type FindOneAndUpdateObject struct {
	ctx     context.Context
	coll    *mongo.Collection
	opts    *options.FindOneAndUpdateOptions
	filter  interface{}
	updater interface{}
}

func (o *FindOneAndUpdateObject) Where(filter interface{}) *FindOneAndUpdateObject {
	o.filter = filter
	return o
}

func (o *FindOneAndUpdateObject) Set(doc interface{}) *FindOneAndUpdateObject {
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

func (o *FindOneAndUpdateObject) Unset(doc interface{}) *FindOneAndUpdateObject {
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

func (o *FindOneAndUpdateObject) Updater(doc interface{}) *FindOneAndUpdateObject {
	o.updater = doc
	return o
}

func (o *FindOneAndUpdateObject) Upsert() *FindOneAndUpdateObject {
	o.opts.SetUpsert(true)
	return o
}

func (o *FindOneAndUpdateObject) Sort(sort interface{}) *FindOneAndUpdateObject {
	o.opts.SetSort(sort)
	return o
}

func (o *FindOneAndUpdateObject) ReturnAfter(returnAfter bool) *FindOneAndUpdateObject {
	returnDoc := options.Before
	if returnAfter {
		returnDoc = options.After
	}
	o.opts.SetReturnDocument(returnDoc)
	return o
}

func (o *FindOneAndUpdateObject) Select(projection interface{}) *FindOneAndUpdateObject {
	o.opts.Projection = projection
	return o
}

func (o *FindOneAndUpdateObject) Do(result interface{}) error {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	rst := o.coll.FindOneAndUpdate(o.ctx, o.filter, o.updater, o.opts)
	return rst.Decode(result)
}

type FindOneAndDeleteObject struct {
	ctx    context.Context
	coll   *mongo.Collection
	opts   *options.FindOneAndDeleteOptions
	filter interface{}
}

func (o *FindOneAndDeleteObject) Where(filter interface{}) *FindOneAndDeleteObject {
	o.filter = filter
	return o
}

func (o *FindOneAndDeleteObject) Sort(sort interface{}) *FindOneAndDeleteObject {
	o.opts.SetSort(sort)
	return o
}

func (o *FindOneAndDeleteObject) Select(projection interface{}) *FindOneAndDeleteObject {
	o.opts.Projection = projection
	return o
}

func (o *FindOneAndDeleteObject) Do(result interface{}) error {
	if o.filter == nil {
		o.filter = bson.D{}
	}
	rst := o.coll.FindOneAndDelete(o.ctx, o.filter, o.opts)
	return rst.Decode(result)
}
