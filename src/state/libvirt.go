package state

import (
	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"
)

func (state *State) SetDomainState(newState libvirt.DomainState, dom libvirt.Domain) {
	// TODO:~ dom.GetState() also returns reason, we may do something with that!
	domState, _, err := dom.GetState()
	if err != nil {
		log.PrintLog(log.WARN, "Error getting domain state! ", err)
	}
	//var isActive bool
	if domState == libvirt.DomainState(8) {
		log.PrintLog(log.ERROR, "Got unexpected domain state! Assuming the domain is shut off.")
	}
	// case libvirt.DOMAIN_NOSTATE:
	// case libvirt.DOMAIN_RUNNING:
	// case libvirt.DOMAIN_BLOCKED:
	// case libvirt.DOMAIN_PAUSED:
	// case libvirt.DOMAIN_SHUTDOWN:
	// case libvirt.DOMAIN_SHUTOFF:
	// case libvirt.DOMAIN_CRASHED:
	// case libvirt.DOMAIN_PMSUSPENDED:
}
