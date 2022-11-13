/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */

package mongu

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexesObject struct {
	ctx     context.Context
	coll    *mongo.Collection
	opts    *options.CreateIndexesOptions
	indexes []mongo.IndexModel
}

func (o *IndexesObject) Keys(keys interface{}) *IndexesObject {
	indexes := append(o.indexes, mongo.IndexModel{Keys: keys, Options: &options.IndexOptions{}})
	o.indexes = indexes
	return o
}

// TTLKey: key should be date specifying the expire time
func (o *IndexesObject) TTLKey(key interface{}, expireSeconds int) *IndexesObject {
	toInt32Seconds := int32(expireSeconds)
	indexes := append(o.indexes, mongo.IndexModel{Keys: key, Options: &options.IndexOptions{ExpireAfterSeconds: &toInt32Seconds}})
	o.indexes = indexes
	return o
}

func (o *IndexesObject) Unique() *IndexesObject {
	indexLen := len(o.indexes)
	if indexLen <= 0 {
		panic("you should call Keys before call Unique on CollHelperIndexes")
	}
	last := o.indexes[indexLen-1]
	last.Options.SetUnique(true)
	return o
}

func (o *IndexesObject) Do() ([]string, error) {
	if o.indexes == nil {
		o.indexes = []mongo.IndexModel{}
	}
	return o.coll.Indexes().CreateMany(o.ctx, o.indexes, o.opts)
}
