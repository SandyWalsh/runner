package sprinter

import (
	"fmt"
)

// AuthZRules defines a set of cgroups and the mapping of a client id to these cgroups when running processes.
type AuthZRules struct {
	ControlGroups  map[string]ControlGroup // [name]ControlGroup
	ClientToCGroup map[string]string       // [client id]ControlGroup name
}

func (a *AuthZRules) getControlGroup(name string) (*ControlGroup, error) {
	if cg, ok := a.ControlGroups[name]; ok {
		return &cg, nil
	}
	return nil, fmt.Errorf("no such control group: %s", name)
}
