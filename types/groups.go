package types

import (
	"github.com/Nerzal/gocloak/v13"
)

type groupWrapper struct {
	group    *gocloak.Group
	children []*groupWrapper
	parent   *groupWrapper
}

func (wrapper *groupWrapper) Group() *gocloak.Group {
	return wrapper.group
}

type GroupWrapper interface {
	Group() *gocloak.Group
	AddChild(*gocloak.Group) GroupWrapper
	SetParent(*gocloak.Group) GroupWrapper
	InheritedAttributes() *map[string][]string
	Children() []GroupWrapper
}

func (wrapper *groupWrapper) AddChild(child *gocloak.Group) GroupWrapper {
	childWrapper := groupWrapper{
		group: child,
	}
	childWrapper.parent = wrapper
	wrapper.children = append(wrapper.children, &childWrapper)
	return &childWrapper
}

func (wrapper *groupWrapper) SetParent(parent *gocloak.Group) GroupWrapper {
	parentWrapper := groupWrapper{
		group: parent,
	}
	wrapper.parent = &parentWrapper
	parentWrapper.children = append(parentWrapper.children, wrapper)
	return &parentWrapper
}

func (wrapper *groupWrapper) InheritedAttributes() *map[string][]string {
	var attrs *map[string][]string
	if wrapper.parent == nil || wrapper.parent.InheritedAttributes() == nil {
		attrMap := make(map[string][]string)
		attrs = &attrMap
	} else {
		attrs = wrapper.parent.InheritedAttributes()
	}
	if wrapper.group.Attributes != nil {
		for key, value := range *wrapper.group.Attributes {
			(*attrs)[key] = value
		}
	}
	return attrs
}

func WrapGroup(group *gocloak.Group) GroupWrapper {
	return &groupWrapper{group: group}
}

func (wrapper *groupWrapper) Children() []GroupWrapper {
	var children []GroupWrapper
	for index := range wrapper.children {
		children = append(children, wrapper.children[index])
	}
	return children
}
