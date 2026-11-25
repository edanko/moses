package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func mergeDocuments(base, override []byte) ([]byte, error) {
	baseNode, err := parseDocument(base)
	if err != nil {
		return nil, fmt.Errorf("%w: embedded defaults: %w", ErrInvalidConfig, err)
	}
	overrideNode, err := parseDocument(override)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	merged := mergeNode(baseNode, overrideNode)
	document := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{merged},
	}
	data, err := yaml.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("%w: merge: %w", ErrInvalidConfig, err)
	}
	return data, nil
}

func parseDocument(data []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		if errors.Is(err, io.EOF) {
			return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, nil
		}
		return nil, err
	}
	if len(document.Content) == 0 {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, nil
	}
	return document.Content[0], nil
}

func mergeNode(base, override *yaml.Node) *yaml.Node {
	if base.Kind != yaml.MappingNode || override.Kind != yaml.MappingNode {
		return override
	}

	merged := &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     base.Tag,
		Style:   base.Style,
		Content: append([]*yaml.Node(nil), base.Content...),
	}

	index := make(map[string]int, len(merged.Content)/2)
	for i := 0; i+1 < len(merged.Content); i += 2 {
		index[merged.Content[i].Value] = i
	}

	for i := 0; i+1 < len(override.Content); i += 2 {
		key, value := override.Content[i], override.Content[i+1]
		if position, ok := index[key.Value]; ok {
			merged.Content[position+1] = mergeNode(merged.Content[position+1], value)
			continue
		}
		merged.Content = append(merged.Content, key, value)
		index[key.Value] = len(merged.Content) - 2
	}

	return merged
}
