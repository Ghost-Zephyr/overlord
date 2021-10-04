package xml

type yesOrNo string

const (
	yes yesOrNo = "yes"
	no  yesOrNo = "no"
)

func (yn yesOrNo) Bool() bool {
	if yn == yes {
		return true
	} else {
		return false
	}
}

type onOrOff string

const (
	on  onOrOff = "on"
	off onOrOff = "off"
)

func (oo onOrOff) Bool() bool {
	if oo == on {
		return true
	} else {
		return false
	}
}
