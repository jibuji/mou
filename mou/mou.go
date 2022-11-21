/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */

/*

## 增删改查的例子：

### 创建mongu.Coll

```
coll := mongu.From(db.Collection(_collName))
```

###增
```
coll.InsertOne(ctx).Do(proj)
```

###删

```
coll.DeleteOne(ctx).Where(bson.M{"_id": ID}).Do()
```

###改

```
coll.UpdateOne(ctx).Where(bson.M{"_id": ID}).Set(task).Do()
```

###查

```
var wfInfo model.Task
err := coll.FindOneByProjectID(ctx).Where(bson.M{"_id": ID}).Do(&wfInfo)
```

*/

package mou

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Coll struct {
	coll *mongo.Collection
}

func From(coll *mongo.Collection) Coll {
	return Coll{coll: coll}
}

func (o Coll) Name() string {
	return o.coll.Name()
}

func (o Coll) Unwrap() *mongo.Collection {
	return o.coll
}

func (o Coll) FindMany(ctx context.Context) *FindManyObject {
	return &FindManyObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Find(),
	}
}

func (o Coll) FindOne(ctx context.Context) *FindOneObject {
	return &FindOneObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.FindOne(),
	}
}

func (o Coll) FindOneAndUpdate(ctx context.Context) *FindOneAndUpdateObject {
	return &FindOneAndUpdateObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.FindOneAndUpdate(),
	}
}

func (o Coll) FindOneAndDelete(ctx context.Context) *FindOneAndDeleteObject {
	return &FindOneAndDeleteObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.FindOneAndDelete(),
	}
}

func (o Coll) InsertMany(ctx context.Context) *InsertManyObject {
	return &InsertManyObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.InsertMany(),
	}
}

func (o Coll) InsertOne(ctx context.Context) *InsertOneObject {
	return &InsertOneObject{
		ctx:  ctx,
		coll: o.coll,
	}
}

func (o Coll) UpdateOne(ctx context.Context) *UpdateOneObject {
	return &UpdateOneObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Update(),
	}
}

func (o Coll) UpdateMany(ctx context.Context) *UpdateManyObject {
	return &UpdateManyObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Update(),
	}
}

func (o Coll) DeleteOne(ctx context.Context) *DeleteOneObject {
	return &DeleteOneObject{
		ctx:  ctx,
		coll: o.coll,
	}
}

func (o Coll) DeleteMany(ctx context.Context) *DeleteManyObject {
	return &DeleteManyObject{
		ctx:  ctx,
		coll: o.coll,
	}
}

func (o Coll) Indexes(ctx context.Context) *IndexesObject {
	return &IndexesObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.CreateIndexes(),
	}
}

func (o Coll) Count(ctx context.Context) *CountObject {
	return &CountObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Count(),
	}
}

func (o Coll) Drop(ctx context.Context) error {
	return o.coll.Drop(ctx)
}

func (o Coll) Aggregate(ctx context.Context) *AggregateObject {
	return &AggregateObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Aggregate().
			SetAllowDiskUse(true).
			//SetMaxAwaitTime 这个参数会引起 cannot set maxtimems on getmore command for a non-awaitdata cursor 这个错误，
			// 但这个错误只在dev环境有，staging环境没有，还需要后续调查，这里先注释掉这个配置，不会造成太大影响
			// SetMaxAwaitTime(time.Minute).
			SetMaxTime(time.Minute),
		pipelines: []bson.D{},
	}
}

func (o Coll) ReplaceOne(ctx context.Context) *ReplaceOneObject {
	return &ReplaceOneObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Replace(),
	}
}

func (o Coll) Distinct(ctx context.Context) *DistinctObject {
	return &DistinctObject{
		ctx:  ctx,
		coll: o.coll,
		opts: options.Distinct(),
	}
}

func (o Coll) WatchChangeStream(ctx context.Context) *ChangeStreamWatchObject {
	return &ChangeStreamWatchObject{
		ctx:      ctx,
		coll:     o.coll,
		pipeline: mongo.Pipeline{},
		opts:     options.ChangeStream(),
	}
}

func (o Coll) Tx(ctx context.Context) *TxObject {
	return &TxObject{
		ctx:    ctx,
		client: o.coll.Database().Client(),
	}
}
