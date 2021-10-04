package libvirt

import (
	"context"

	"bitsnthings.dev/overlord/src/libvirt/xml"
	"bitsnthings.dev/overlord/src/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"libvirt.org/go/libvirt"
)

type DomainDocument struct {
	ID     primitive.ObjectID `bson:"_id"`
	URI    string             `bson:"url"`
	UUID   string             `bson:"uuid"`
	Active bool               `bson:"active"`
	XML    string             `bson:"xml"`
	domXML xml.DomainXML
}

func domainError(
	err error,
	str string,
	domain DomainDocument,
) bool {
	if err != nil {
		log.PrintLog(
			log.DEBUG,
			str+" domain with uuid: \"%s\" on node with uri: \"%s\" in a database change handler! %s",
			domain.UUID, domain.URI, err,
		)
		return true
	} else {
		return false
	}
}

func (virt *Libvirt) DomainChangeHandler(stream *mongo.ChangeStream) {
	for stream.Next(context.TODO()) {
		go func() {
			updatedDomain := struct {
				FullDocument DomainDocument `bson:"fullDocument"`
			}{}
			if err := stream.Decode(&updatedDomain); err != nil {
				log.PrintLog(log.ERROR, "Error decoding mongodb change stream! %s", err)
			}
			domainError(updatedDomain.FullDocument.domXML.Unmarshal(
				[]byte(updatedDomain.FullDocument.XML),
			), "Error decoding xml for", updatedDomain.FullDocument)
			//!~ This runs without error while not actually decoding anything!
			// log.PrintLog(log.DEBUG, "Decoded domain change: %d", updatedDomain)
			domain := updatedDomain.FullDocument
			if host, exists := virt.Hosts[domain.URI]; exists {
				dom, err := host.LookupDomainByUUIDString(domain.UUID)
				if !domainError(err, "Error getting", domain) {
					name, err := dom.GetName()
					domainError(err, "Error getting name of", domain)
					if name != domain.domXML.Name {
						log.PrintLog(
							log.INFO,
							"Detected name change for domain with uuid: \"%s\" on node with uri: \"%s\"!",
							domain.UUID, domain.URI,
						)
						log.PrintLog(log.DEBUG, "This means we have to undefine the libvirt xml and redefine it!")
						dom.Undefine()
					}
				}
				dom, err = host.DomainDefineXML(domain.XML)
				if domainError(err, "Error updating", domain) {
					return
				}
				setDomainActive(dom, domain)
				log.PrintLog(
					log.INFO,
					"Updated domain with uuid: \"%s\" on node with uri: \"%s\"!",
					domain.UUID, domain.URI)
			} else {
				if _, exists := virt.ROHosts[domain.URI]; exists {
					log.PrintLog(
						log.WARN,
						"Got db update for domain with uuid: \"%s\" on read only host with uri: \"%s\"!",
						domain.UUID, domain.URI,
					)
				} else {
					log.PrintLog(
						log.ERROR,
						"Got db update for domain with uuid: \"%s\" on host with uri: \"%s\", but is not in internal host map!",
						domain.UUID, domain.URI,
					)
				}
			}
		}()
	}
}

func setDomainActive(dom *libvirt.Domain, domain DomainDocument) {
	active, err := dom.IsActive()
	if domainError(err, "Error getting active status of", domain) {
		return
	}
	if !active && domain.Active {
		domainError(dom.Resume(), "Error resuming", domain)
	}
	if active && !domain.Active {
		domainError(dom.Suspend(), "Error suspending", domain)
	}
}
