package mou

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChangeStreamWatchObject struct {
	ctx         context.Context
	coll        *mongo.Collection
	opts        *options.ChangeStreamOptions
	pipeline    mongo.Pipeline
	resumeToken interface{}
}

func (o *ChangeStreamWatchObject) FullDocumentOnUpdate(full bool) *ChangeStreamWatchObject {
	if full {
		o.opts.SetFullDocument(options.UpdateLookup)
	}
	return o
}

func (o *ChangeStreamWatchObject) Pipeline(pipeline mongo.Pipeline) *ChangeStreamWatchObject {
	if pipeline != nil {
		o.pipeline = pipeline
	}
	return o
}

func (o *ChangeStreamWatchObject) ResumeAfter(resumeToken interface{}) *ChangeStreamWatchObject {
	o.resumeToken = resumeToken
	return o
}

type WrUpdateDescription struct {
	UpdatedFields   bson.M   `bson:"updatedFields"`
	RemovedFields   []string `bson:"removedFields"`
	TruncatedArrays []bson.M `bson:"truncatedArrays"`
}

type WrDocumentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type WrNamespace struct {
	DB   string `bson:"db"`
	Coll string `bson:"coll"`
}

// type WatchResult2 struct {
// 	OperationType     string
// 	Document          interface{}
// 	ResumeToken       bson.M
// 	DocumentKey       WrDocumentKey
// 	UpdateDescription WrUpdateDescription
// }

// // Deprecated:
// //out will be automatically closed when no value available or error happened
// func (o *ChangeStreamWatchObject) Do2() (<-chan WatchResult2, error) {
// 	if o.resumeToken != nil {
// 		o.opts.SetResumeAfter(o.resumeToken)
// 		//.SetFullDocument(options.UpdateLookup)
// 	}
// 	// Watch the collection
// 	cs, err := o.coll.Watch(o.ctx, mongo.Pipeline{}, o.opts)
// 	if err != nil {
// 		fmt.Printf("Watch change stream for %s error: %v \n we will try again without resumeToken", o.coll.Name(), err)
// 		cs, err = o.coll.Watch(o.ctx, mongo.Pipeline{}, o.opts.SetResumeAfter(nil))
// 		if err != nil {
// 			//still error, return out
// 			return nil, err
// 		}
// 	}
// 	if o.docSelector == nil {
// 		o.docSelector = struct{}{}
// 	}
// 	// var ID changeID
// 	eventType := reflect.StructOf([]reflect.StructField{
// 		{
// 			Name: "ID",
// 			Type: reflect.TypeOf(bson.M{}),
// 			Tag:  `bson:"_id"`,
// 		},
// 		{
// 			Name: "OperationType",
// 			Type: reflect.TypeOf(""),
// 			Tag:  `bson:"operationType"`,
// 		},
// 		{
// 			Name: "FullDocument",
// 			Type: reflect.TypeOf(o.docSelector),
// 			Tag:  `bson:"fullDocument"`,
// 		},
// 		{
// 			Name: "DocumentKey",
// 			Type: reflect.TypeOf(bson.M{}),
// 			Tag:  `bson:"documentKey"`,
// 		},
// 		{
// 			Name: "UpdateDescription",
// 			Type: reflect.TypeOf(WrUpdateDescription{}),
// 			Tag:  `bson:"updateDescription"`,
// 		},
// 	})
// 	out := make(chan WatchResult2, 100)
// 	// Whenever there is a new change event, decode the change event and print some information about it
// 	go func() {
// 		defer close(out)
// 		for cs.Next(o.ctx) {
// 			pv := reflect.New(eventType)
// 			ev := pv.Interface()
// 			err := cs.Decode(ev)
// 			if err != nil {
// 				log.Printf("%v", err)
// 				return
// 			}
// 			pve := pv.Elem()
// 			wr := WatchResult2{}
// 			wr.ResumeToken = pve.FieldByName("ID").Interface().(bson.M)
// 			wr.Document = pve.FieldByName("FullDocument").Interface()
// 			wr.OperationType = pve.FieldByName("OperationType").String()
// 			wr.DocumentKey = pve.FieldByName("DocumentKey").Interface().(bson.M)
// 			wr.UpdateDescription = pve.FieldByName("UpdateDescription").Interface().(WrUpdateDescription)
// 			out <- wr
// 		}
// 	}()
// 	return out, nil
// }

type WatchResult struct {
	ID                string              `bson:"-"`
	ClusterTime       primitive.Timestamp `bson:"clusterTime"`
	Namespace         WrNamespace         `bson:"ns"`
	ResumeToken       bson.M              `bson:"_id"`
	OperationType     string              `bson:"operationType"`
	FullDocument      bson.M              `bson:"fullDocument"`
	DocumentKey       WrDocumentKey       `bson:"documentKey"`
	UpdateDescription WrUpdateDescription `bson:"updateDescription"`
}

// out will be automatically closed when no value available or error happened
func (o *ChangeStreamWatchObject) Do() (<-chan WatchResult, error) {
	if o.resumeToken != nil {
		o.opts.SetResumeAfter(o.resumeToken)
	}
	// Watch the collection
	cs, err := o.coll.Watch(o.ctx, o.pipeline, o.opts)
	if err != nil {
		fmt.Printf("Watch change stream for %s error: %v \n we will try again without resumeToken", o.coll.Name(), err)
		cs, err = o.coll.Watch(o.ctx, mongo.Pipeline{}, o.opts.SetResumeAfter(nil))
		if err != nil {
			//still error, return out
			return nil, err
		}
	}

	out := make(chan WatchResult, 100)
	// Whenever there is a new change event, decode the change event and print some information about it
	go func() {
		defer close(out)
		for cs.Next(o.ctx) {
			wr := WatchResult{}
			err := cs.Decode(&wr)
			if err != nil {
				// log.Printf("%v", err)
				return
			}
			// log.Printf("wr.ResumeToken=%v, wr.ResumeToken['_id']=%v", wr.ResumeToken, wr.ResumeToken["_id"])
			wr.ID = wr.ResumeToken["_data"].(string)
			out <- wr
		}
	}()
	return out, nil
}
