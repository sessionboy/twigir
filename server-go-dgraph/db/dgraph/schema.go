package dgraph

import (
	"context"
	dschema "server/db/dgraph/schema"

	"github.com/dgraph-io/dgo/v210/protos/api"
	log "github.com/sirupsen/logrus"
)

func InitSchema() (err error) {
	dg, close, err := NewClient()
	if err != nil {
		log.Errorf("cannot init dgraph schema: %v", err)
		// TODO retry after few seconds
		return err
	}
	defer close()

	schema := dschema.UserSchema + dschema.StatusSchema + dschema.NotificationSchema + dschema.MessageSchema
	log.Infof("schema: %s", schema)
	ctx := context.Background()
	err = dg.Alter(ctx, &api.Operation{
		Schema: string(schema),
	})
	if err != nil {
		log.Errorf("init schema fail: %v", err)
		return err
	}

	log.Infof("schema initialized successfull!")
	return nil
}
