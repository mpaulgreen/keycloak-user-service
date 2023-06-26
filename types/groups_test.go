package types

import (
	"github.com/Nerzal/gocloak/v13"
	"testing"
)

func TestWrapGroup(t *testing.T) {
	parent := WrapGroup(&gocloak.Group{})
	if _, ok := parent.(*groupWrapper); !ok {
		t.Fatalf("Did not create wrapper as expected")
	}
}

func TestAddingChildren(t *testing.T) {
	parentName := "Parent"
	firstChildName := "Child 1"
	secondChildName := "Child 2"
	parent := WrapGroup(&gocloak.Group{Name: &parentName})
	parent.AddChild(&gocloak.Group{Name: &firstChildName})
	parent.AddChild(&gocloak.Group{Name: &secondChildName})
	if *parent.(*groupWrapper).group.Name != parentName {
		t.Errorf("Unexpected parent name")
	}
	if len(parent.(*groupWrapper).children) != 2 {
		t.Errorf("Expect two child groups")
	}
	if *parent.(*groupWrapper).children[0].group.Name != firstChildName {
		t.Errorf("Expect first child entry to be " + firstChildName)
	}
	if *parent.(*groupWrapper).children[1].group.Name != secondChildName {
		t.Errorf("Expect second child entry to be " + secondChildName)
	}
}

func TestAttributeRetention(t *testing.T) {
	parentAttrs := make(map[string][]string)
	parentAttrs["orgId"] = []string{"1010101"}
	parentAttrs["approved"] = []string{"false"}
	parent := WrapGroup(&gocloak.Group{Attributes: &parentAttrs})

	node1Attrs := make(map[string][]string)
	node1Attrs["approved"] = []string{"true"}
	child1 := parent.AddChild(&gocloak.Group{Attributes: &node1Attrs})

	node2Attrs := make(map[string][]string)
	node2Attrs["custom"] = []string{"has"}
	child2 := parent.AddChild(&gocloak.Group{Attributes: &node2Attrs})

	if (*(parent.(*groupWrapper).group.Attributes))["orgId"][0] != "1010101" {
		t.Errorf("Expected orgId to be set to %d", 1010101)
	}
	if (*(parent.(*groupWrapper).group.Attributes))["approved"][0] != "false" {
		t.Errorf("Expected approved to be set to %s", "false")
	}
	if (*(child1.(*groupWrapper).group.Attributes))["approved"][0] != "true" {
		t.Errorf("Expected approved to be set to %s", "true")
	}
	if (*(child2.(*groupWrapper).group.Attributes))["custom"][0] != "has" {
		t.Errorf("Expected custom to be set to %s", "has")
	}
}

func TestAttributeInhertiance(t *testing.T) {
	parentAttrs := make(map[string][]string)
	parentAttrs["orgId"] = []string{"1010101"}
	parentAttrs["approved"] = []string{"false"}
	parent := WrapGroup(&gocloak.Group{Attributes: &parentAttrs})

	node1Attrs := make(map[string][]string)
	node1Attrs["approved"] = []string{"true"}
	child1 := parent.AddChild(&gocloak.Group{Attributes: &node1Attrs})
	effective1 := *child1.InheritedAttributes()

	node2Attrs := make(map[string][]string)
	node2Attrs["custom"] = []string{"has"}
	child2 := parent.AddChild(&gocloak.Group{Attributes: &node2Attrs})
	effective2 := *child2.InheritedAttributes()

	if effective1["orgId"][0] != "1010101" {
		t.Errorf("Expected orgId to be set to %d", 1010101)
	}
	if effective1["approved"][0] != "true" {
		t.Errorf("Expected approved to be set to %s", "true")
	}
	if effective2["orgId"][0] != "1010101" {
		t.Errorf("Expected orgId to be set to %d", 1010101)
	}
	if effective2["custom"][0] != "has" {
		t.Errorf("Expected approved to be set to %s", "has")
	}
}
