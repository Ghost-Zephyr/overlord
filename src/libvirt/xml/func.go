package xml

import "encoding/xml"

func (domXml *DomainXML) Unmarshal(newXml []byte) error {
	return xml.Unmarshal(newXml, &domXml)
}
