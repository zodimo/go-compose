package tree

import (
	"net/url"
	"strings"
)

type TreeId struct {
	encodedLineage []string
	decodedLabel   string
	labelEndoder   LabelEncoder
}

func (id TreeId) Equals(oth TreeId) bool {
	return id.String() == oth.String()
}

func (id TreeId) String() string {
	return strings.Join(append(id.encodedLineage, id.labelEndoder.Encode(id.decodedLabel)), "/")
}

func (id TreeId) Label() string {
	return id.labelEndoder.Decode(id.decodedLabel)
}

func NewTreeNodeId(label string, labelEndoder ...LabelEncoder) TreeId {

	var localLabelEncoder LabelEncoder
	if len(labelEndoder) >= 1 {
		localLabelEncoder = labelEndoder[0]
	}

	if localLabelEncoder == nil {
		localLabelEncoder = PassthruEncoder
	}

	return TreeId{
		encodedLineage: []string{},
		decodedLabel:   label,
		labelEndoder:   localLabelEncoder,
	}
}

func NewChildTreeId(label string, parentId TreeId) TreeId {
	lineage := append(parentId.encodedLineage, parentId.labelEndoder.Encode(parentId.decodedLabel))
	return TreeId{
		encodedLineage: lineage,
		decodedLabel:   label,
		labelEndoder:   parentId.labelEndoder,
	}
}

type LabelEncoder interface {
	Encode(string) string
	Decode(string) string
}

var _ LabelEncoder = (*labelEncoderImpl)(nil)

type labelEncoderImpl struct {
	encode func(string) string
	decode func(string) string
}

func NewLabelEncoder(encode func(string) string, decode func(string) string) LabelEncoder {
	if encode == nil {
		panic("The encode function cannot be nil")
	}
	if decode == nil {
		panic("The decode function cannot be nil")
	}
	return labelEncoderImpl{
		encode: encode,
		decode: decode,
	}
}

func FromString(idString string, labelEndoder ...LabelEncoder) TreeId {
	var localLabelEncoder LabelEncoder
	if len(labelEndoder) >= 1 {
		localLabelEncoder = labelEndoder[0]
	}

	if localLabelEncoder == nil {
		localLabelEncoder = PassthruEncoder
	}

	parts := strings.Split(idString, "/")

	if len(parts) == 1 {
		return NewTreeNodeId(localLabelEncoder.Decode(parts[0]), localLabelEncoder)
	}

	return TreeId{
		encodedLineage: parts[:len(parts)-1],
		decodedLabel:   localLabelEncoder.Decode(parts[len(parts)-1]),
		labelEndoder:   localLabelEncoder,
	}
}

func (le labelEncoderImpl) Encode(label string) string {
	return le.encode(label)
}
func (le labelEncoderImpl) Decode(label string) string {
	return le.decode(label)
}

var PassthruEncoder = NewLabelEncoder(
	func(x string) string { return x },
	func(x string) string { return x },
)

var QueryEscapeEncoder = NewLabelEncoder(
	url.QueryEscape,
	func(x string) string {
		res, err := url.QueryUnescape(x)
		if err != nil {
			return x
		}
		return res
	},
)
