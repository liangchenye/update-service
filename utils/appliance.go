package utils

import (
	"fmt"
)

type Appliance struct {
	Proto      string
	Version    string
	Namespace  string
	Repository string
	//only for 'app'
	OS string
	//only for 'app'
	Arch string
	Name string
	Tag  string
}

func (a *Appliance) FullName() string {
	var val string
	var tag string

	if a.Tag != "" {
		tag = a.Tag
	} else {
		tag = "latest"
	}

	switch a.Proto {
	case "app":
		val = fmt.Sprintf("%s-%s-%s:%s", a.OS, a.Arch, a.Name, tag)
	default:
		val = fmt.Sprintf("%s:%s", a.Name, tag)
	}

	return val
}
