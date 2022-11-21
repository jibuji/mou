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

type AggregateObject struct {
	ctx       context.Context
	coll      *mongo.Collection
	opts      *options.AggregateOptions
	pipelines mongo.Pipeline
}

func (o *AggregateObject) DoCount() (int64, error) {
	o.pipelines = append(o.pipelines, bson.D{{Key: "$count", Value: "count"}})
	var result []struct {
		Count int64 `bson:"count"`
	}
	err := o.Do(&result)
	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}
	return count, err
}

// func (o *AggregateObject) Group(stage interface{}) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$group", Value: stage}})
// 	return o
// }

// func (o *AggregateObject) GraphLookup(stage bson.D) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$graphLookup", Value: stage}})
// 	return o
// }

// func (o *AggregateObject) Skip(i int64) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$skip", Value: i}})
// 	return o
// }
// func (o *AggregateObject) Limit(i int64) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$limit", Value: i}})
// 	return o
// }

// func (o *AggregateObject) Sort(stage interface{}) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$sort", Value: stage}})
// 	return o
// }

// func (o *AggregateObject) AddFields(stage interface{}) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$addFields", Value: stage}})
// 	return o
// }

// func (o *AggregateObject) UnionWith(coll string, pipeline []interface{}) *AggregateObject {
// 	o.pipelines = append(o.pipelines, bson.D{{Key: "$unionWith", Value: bson.D{{Key: "coll", Value: coll}, {Key: "pipeline", Value: pipeline}}}})
// 	return o
// }

func (o *AggregateObject) Pipelines(p mongo.Pipeline) *AggregateObject {
	o.pipelines = p
	return o
}

// func (o *AggregateObject) Clone() *AggregateObject {
// 	pipelines := append([]bson.D{}, o.pipelines...)
// 	return &AggregateObject{ctx: o.ctx, coll: o.coll, opts: o.opts, pipelines: pipelines}
// }

func (o *AggregateObject) Do(results interface{}) error {
	if o.pipelines == nil {
		o.pipelines = []bson.D{}
	}
	cursor, err := o.coll.Aggregate(o.ctx, o.pipelines, o.opts)
	if err != nil {
		return err
	}
	defer cursor.Close(o.ctx)
	return cursor.All(o.ctx, results)
}
