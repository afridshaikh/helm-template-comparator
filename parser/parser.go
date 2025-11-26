package parser

import (
	"fmt"
	kindregistry "helm_template_comparator/kind_registry"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type Resource struct {
	Key  string
	Kind string
	Obj  interface{}
}

func splitMultiYAML(data []byte) [][]byte {
	docs := strings.Split(string(data), "---")
	var result [][]byte
	for _, d := range docs {
		trimmed := strings.TrimSpace(d)
		if trimmed != "" {
			result = append(result, []byte(trimmed))
		}
	}
	return result
}

func parseResource(doc []byte) (*Resource, error) {
	var meta struct {
		Kind string `yaml:"kind"`
		Metadata struct {
			Name string `yaml:"name"`
		}`yaml:"metadata"`
	}
	if err := yaml.Unmarshal(doc, &meta); err != nil {
		return nil, err
	}
	if meta.Metadata.Name == "" {
		return nil, fmt.Errorf("resource missing name")
	}

	kindInfo, ok := kindregistry.KindRegistry[meta.Kind]
	if !ok {
		return nil, fmt.Errorf("unsupported kind: %s", meta.Kind)
	}

	obj := kindInfo.Type()
	if err := yaml.Unmarshal(doc, obj); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s", meta.Metadata.Name, meta.Kind)
	return &Resource{Key: key, Obj: obj, Kind: meta.Kind}, nil
}

func ReadHelmTemplate(path string) (map[string]*Resource, error) {
	objs := make(map[string]*Resource)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for _, doc := range splitMultiYAML(data) {
		res, err := parseResource(doc)
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", path, err)
			continue
		}
		objs[res.Key] = res
	}

	return objs, nil
}
