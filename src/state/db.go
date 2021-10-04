package state

import (
	"context"

	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (state *State) setupDB() {
	state.upsertOpts.SetUpsert(true)
	state.streamOpts.SetBatchSize(8)
	state.streamOpts.SetFullDocument("updateLookup")
	if state.Config.MongoDbUri == "" {
		log.PrintLog(log.INFO, "No MongoDB uri in config, using in memory database!")
	}
	if state.Config.InMemoryDB {
		log.PrintLog(log.INFO, "InMemoryDB set in config, using in memory database!")
	}
	if state.Config.MongoDbUri == "" || state.Config.InMemoryDB {
		log.PrintLog(log.WARN, "MongoDB collection watchers is not supported with in memory database!")
		state.Config.InMemoryDB = true
		state.Config.MongoDbUri = ""
		state.DB.domains = &mongo.Collection{}
		return
	}
	clientOptions := options.Client().ApplyURI(state.Config.MongoDbUri)
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	err := client.Ping(context.TODO(), nil)
	// No idea why, but error only gets returned with ping, not connect,
	// if the server refuses the connection.
	if err != nil {
		log.PrintLog(log.FATAL, "Error connecting to database wtih connection string \"%s\"! %s",
			state.Config.MongoDbUri, err)
	}
	dbName := state.Config.MongoDbName
	if dbName == "" {
		dbName = "overlord"
	}
	db := client.Database(dbName)
	state.DB.domains = db.Collection("domains")
	state.setDBWatchers()
}

func (state *State) setDBWatchers() {
	state.startWatchStreamHandler(
		*state.DB.domains,
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
	for uuid, dom := range domMap {
		xml, err := dom.GetXMLDesc(state.Libvirt.GetXMLFlags(uri))
		if err != nil {
			log.PrintLog(
				log.ERROR,
				"Error getting XML of domain with uuid: \"%s\" on node with uri: \"%s\"! %s",
				uuid, uri, err)
		}
		state.DB.domains.UpdateOne(
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
