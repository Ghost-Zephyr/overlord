package state

import (
	"bitsnthings.dev/overlord/src/api"
	virt "bitsnthings.dev/overlord/src/libvirt"
	"bitsnthings.dev/overlord/src/state/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	domains *mongo.Collection
}

type State struct {
	Libvirt       virt.Libvirt
	Config        conf.Config
	DB            DB
	upsertOpts    *options.UpdateOptions
	streamOpts    *options.ChangeStreamOptions
	changeStreams []*mongo.ChangeStream
	API           api.API
}
