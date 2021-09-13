package libvirt

import (
	"context"

	"bitsnthings.dev/overlord/src/log"
	"go.mongodb.org/mongo-driver/mongo"
	"libvirt.org/go/libvirt"
)

type domainDocument struct {
	Active bool
	Uuid   string
	Uri    string
	Xml    string
}

func domainError(
	err error,
	str string,
	domain domainDocument,
	a ...interface{},
) bool {
	if err != nil {
		log.PrintLog(
			log.ERROR,
			str+" domain with uuid: \"%s\" on node with uri: \"%s\"!",
			a, domain.Uuid, domain.Uri,
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
				FullDocument domainDocument `bson:"fullDocument"`
			}{}
			if err := stream.Decode(&updatedDomain); err != nil {
				log.PrintLog(log.ERROR, "Error decoding mongodb change stream! %s", err)
			}
			domain := updatedDomain.FullDocument
			dom, err := virt.Hosts[domain.Uri].DomainDefineXML(domain.Xml)
			if domainError(err, "Error updating", domain) {
				return
			}
			setDomainActive(dom, domain)
			log.PrintLog(
				log.INFO,
				"Updated domain with uuid: \"%s\" on node with uri: \"%s\"!",
				domain.Uuid, domain.Uri)
		}()
	}
}

func setDomainActive(dom *libvirt.Domain, domain domainDocument) {
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
