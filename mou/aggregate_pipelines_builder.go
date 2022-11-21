/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-10-13 08:31:06
 */

package mou

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PipelineBuilderType mongo.Pipeline
type builder = PipelineBuilderType

func PipelineBuilder() builder {
	return builder{}
}

func PipelineJsonBuilder(pipelineJSON []byte) (builder builder, err error) {
	err = bson.UnmarshalExtJSON([]byte(pipelineJSON), true, &builder)
	return
}

func (o builder) Match(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$match", Value: stage}})
	return o
}

func (o builder) Lookup(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$lookup", Value: stage}})
	return o
}

func (o builder) Project(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$project", Value: stage}})
	return o
}

func (o builder) Unwind(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$unwind", Value: stage}})
	return o
}

func (o builder) ReplaceRoot(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$replaceRoot", Value: stage}})
	return o
}

func (o builder) Unset(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$unset", Value: stage}})
	return o
}

func (o builder) Group(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$group", Value: stage}})
	return o
}

func (o builder) GraphLookup(stage bson.D) builder {
	o = append(o, bson.D{{Key: "$graphLookup", Value: stage}})
	return o
}

func (o builder) Skip(i int64) builder {
	o = append(o, bson.D{{Key: "$skip", Value: i}})
	return o
}

func (o builder) Limit(i int64) builder {
	o = append(o, bson.D{{Key: "$limit", Value: i}})
	return o
}

func (o builder) Sort(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$sort", Value: stage}})
	return o
}

func (o builder) AddFields(stage interface{}) builder {
	o = append(o, bson.D{{Key: "$addFields", Value: stage}})
	return o
}

func (o builder) UnionWith(coll string, pipeline []interface{}) builder {
	o = append(o, bson.D{{Key: "$unionWith", Value: bson.D{{Key: "coll", Value: coll}, {Key: "pipeline", Value: pipeline}}}})
	return o
}

func (o builder) copy() builder {
	other := make(builder, len(o))
	copy(other, o)
	return other
}

func (o builder) Append(builderOrPipeline interface{}) builder {
	switch value := builderOrPipeline.(type) {
	case PipelineBuilderType:
		o = append(o, value...)
	case mongo.Pipeline:
		o = append(o, value...)
	case []bson.D:
		o = append(o, value...)
	default:
		err := fmt.Errorf("unsupported type:%v", reflect.TypeOf(builderOrPipeline))
		panic(err)
	}
	return o
}

func (o builder) ToPipeline() mongo.Pipeline {
	return mongo.Pipeline(o.copy())
}
