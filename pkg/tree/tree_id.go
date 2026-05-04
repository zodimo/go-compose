package tree

import "strings"

type TreeId struct {
	lineage   []string
	label     string
	separator string
}

func (id TreeId) String() string {
	separator := id.separator
	if separator == "" {
		separator = "/"
	}
	return strings.Join(append(id.lineage, id.label), separator)
}

func (id TreeId) Label() string {
	return id.label
}

func NewTreeNodeId(label string) TreeId {
	return TreeId{
		lineage: []string{},
		label:   label,
	}
}

func NewChildTreeId(label string, parentId TreeId) TreeId {
	lineage := append(parentId.lineage, parentId.label)
	return TreeId{
		lineage: lineage,
		label:   label,
	}
}
