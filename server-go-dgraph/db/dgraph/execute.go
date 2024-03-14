package dgraph

import (
	"context"
	"encoding/json"
	common "server/models/common"

	dgo "github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	log "github.com/sirupsen/logrus"
)

func Query(ctx context.Context, txn *dgo.Txn, q string) (r common.Object, err error) {
	resp, err := txn.Query(ctx, q)
	if err != nil {
		return
	}
	var result = common.Object{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return
	}
	r = result
	return
}

func QueryWithVars(ctx context.Context, txn *dgo.Txn, q string, vars common.StrObject) (r common.Object, err error) {
	resp, err := txn.QueryWithVars(ctx, q, vars)
	// fmt.Println("resp", resp)
	// fmt.Println("errr", err)
	if err != nil {
		return
	}
	var result = common.Object{}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return
	}
	r = result
	return
}

func Mutate(ctx context.Context, txn *dgo.Txn, m interface{}) (r common.Object, err error) {
	pb, err := json.Marshal(m)
	if err != nil {
		log.Println("json.Marshal err:: ", err)
		return
	}
	mu := &api.Mutation{
		SetJson: pb,
	}
	resp, err := txn.Mutate(ctx, mu)

	if err != nil {
		log.Println("txn.Mutate resp:: ", resp)
		log.Println("txn.Mutate err:: ", err)
		return
	}
	var result = common.Object{}
	result["uids"] = resp.GetUids()
	if len(resp.Json) > 0 {
		err = json.Unmarshal(resp.Json, &result)
		if err != nil {
			log.Println("json.Unmarshal err:: ", err)
			return
		}
	}
	return result, nil
}

func MutateWithNquads(ctx context.Context, txn *dgo.Txn, nquads []byte) (r common.Object, err error) {
	mu := &api.Mutation{
		SetNquads: nquads,
	}
	resp, err := txn.Mutate(ctx, mu)
	if err != nil {
		return
	}
	var result = common.Object{}
	result["uids"] = resp.GetUids()
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return
	}
	r = result
	return
}

func DeleteWithNquads(ctx context.Context, txn *dgo.Txn, nquads []byte) (r common.Object, err error) {
	mu := &api.Mutation{
		DelNquads: nquads,
	}
	resp, err := txn.Mutate(ctx, mu)
	if err != nil {
		return
	}
	var result = common.Object{}
	result["uids"] = resp.GetUids()
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return
	}
	r = result
	return
}

func GetItem(o common.Object, key string) interface{} {
	items := o[key].([]interface{})
	if len(items) == 0 {
		return nil
	}
	return items[0]
}

func GetList(o common.Object, key string) interface{} {
	items := o[key].([]interface{})
	if len(items) == 0 {
		return nil
	}
	return items
}

func GetSubList(o common.Object, key string, sub_key string) interface{} {
	var list interface{}
	items := o[key].([]interface{})
	if len(items) == 0 {
		list = make([]interface{}, 0)
	} else {
		item := items[0].(map[string]interface{})
		if item[sub_key] == nil {
			list = make([]interface{}, 0)
		} else {
			list = item[sub_key]
		}
	}
	return list
}

func GetUids(o common.Object) map[string]string {
	if o["uids"] == nil {
		return nil
	}
	return o["uids"].(map[string]string)
}

func GetUid(o common.Object, key string) string {
	if o["uids"] == nil {
		return ""
	} else {
		uids := o["uids"].(map[string]string)
		return uids[key]
	}
}
