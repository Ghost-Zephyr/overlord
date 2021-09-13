package state

import (
	"context"

	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (state *State) setDBWatchers() {
	state.startWatchStreamHandler(
		*state.MongoDB.Collection("domains"),
		mongo.Pipeline{},
		state.Libvirt.DomainChangeHandler)
}

func (state *State) startWatchStreamHandler(
	collection mongo.Collection,
	pipeline mongo.Pipeline,
	handler func(*mongo.ChangeStream),
) {
	stream, err := collection.Watch(
		context.TODO(),
		pipeline,
		state.streamOpts)
	if err != nil {
		log.PrintLog(
			log.ERROR,
			"Error starting mongodb watch stream! %s",
			err)
		stream.Close(context.TODO())
		return
	}
	state.changeStreams = append(state.changeStreams, stream)
	go handler(stream)
}

func (state *State) fetchDBState() {

}

func (state *State) pushDBState() {
	for uri, domMap := range state.Libvirt.Domains.Active {
		setDomainStatesInDB(
			domMap, state, true, uri, state.upsertOpts)
	}
	for uri, domMap := range state.Libvirt.Domains.Inactive {
		setDomainStatesInDB(
			domMap, state, false, uri, state.upsertOpts)
	}
}

func setDomainStatesInDB(
	domMap map[string]libvirt.Domain,
	state *State,
	domainState bool,
	uri string,
	opts *options.UpdateOptions,
) {
	var flags libvirt.DomainXMLFlags
	if !state.Libvirt.IsReadOnly(uri) {
		flags = libvirt.DOMAIN_XML_INACTIVE + libvirt.DOMAIN_XML_INACTIVE
	}
	for uuid, dom := range domMap {
		xml, err := dom.GetXMLDesc(flags)
		if err != nil {
			log.PrintLog(
				log.ERROR,
				"Error getting XML of domain with uuid: \"%s\" on node with uri: \"%s\"! %s",
				uuid, uri, err)
		}
		state.MongoDB.Collection("domains").UpdateOne(
			context.TODO(),
			bson.D{
				{Key: "uri", Value: uri},
				{Key: "uuid", Value: uuid}},
			bson.D{
				{Key: "$set", Value: bson.M{
					"active": domainState,
					"xml":    xml,
				}}},
			opts,
		)
	}
}
