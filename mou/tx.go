/*
 * @Descripttion:
 * @version:
 * @Author: pengfei
 * @Date: 2021-07-01 08:31:06
 */
package mou

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TxFn = func(sessCtx mongo.SessionContext) (interface{}, error)
type TxObject struct {
	ctx    context.Context
	client *mongo.Client
}

func (o *TxObject) Run(fn TxFn) error {
	return RunInCtxTransaction(o.ctx, o.client, fn)
}

func RunInCtxTransaction(ctx context.Context, client *mongo.Client, callback TxFn) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, callback)
	return err
}

func RunInTransaction(client *mongo.Client, callback TxFn) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, callback)
	return err
}
