package libvirt

import (
	"encoding/xml"

	virtxml "bitsnthings.dev/overlord/src/libvirt/xml"
	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"
)

func (virt *Libvirt) GetAllDomains() {
	virt.Domains.Active = make(map[string]map[string]libvirt.Domain)
	virt.Domains.Inactive = make(map[string]map[string]libvirt.Domain)
	for uri, conn := range virt.Hosts {
		virt.Domains.Active[uri] = virt.getDomainsByFlag(uri, conn, libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		virt.Domains.Inactive[uri] = virt.getDomainsByFlag(uri, conn, libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	}
}

func (virt *Libvirt) getDomainsByFlag(uri string, conn *libvirt.Connect, flag libvirt.ConnectListAllDomainsFlags) map[string]libvirt.Domain {
	domMap := make(map[string]libvirt.Domain)
	if _, inMap := virt.DomainXMLs[uri]; !inMap {
		virt.DomainXMLs[uri] = make(map[string]virtxml.DomainXML)
	}
	doms, err := conn.ListAllDomains(flag)
	if err != nil {
		log.PrintLog(
			log.ERROR,
			"Error fetching domain list from node with connection uri: \"%s\"! %s",
			uri, err)
	}
	for _, dom := range doms {
		uuid, err := dom.GetUUIDString()
		if err != nil {
			log.PrintLog(
				log.ERROR,
				"Error getting uuid of domain on node with uri: \"%s\"! %s",
				uri, err)
		} else {
			domMap[uuid] = dom
			go func() {
				if stateXML, inMap := virt.DomainXMLs[uri][uuid]; inMap {
					if activeXMLstr, err := dom.GetXMLDesc(virt.GetXMLFlags(uri)); err != nil {
						var activeXML virtxml.DomainXML
						xml.Unmarshal([]byte(activeXMLstr), &activeXML)
						log.PrintLog(log.DEBUG, "Parsed domain XML: %d", activeXML)
						if stateXML != activeXML {
							virt.DomainXMLs[uri][uuid] = activeXML
						}
					}
				}
			}()
		}
	}
	return domMap
}
