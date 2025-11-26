package main

import (
	"fmt"
	kindregistry "helm_template_comparator/kind_registry"
	"helm_template_comparator/parser"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
)

const (
	PROD = "/Users/afridshaikh/projects/golang/junk/prod/"
	DR   = "/Users/afridshaikh/projects/golang/junk/prod-dr/"
)

func compareResources(res1, res2 *parser.Resource) {
	if kindInfo, ok := kindregistry.KindRegistry[res1.Kind]; ok {
		opts := []cmp.Option{}
		if kindInfo.CompareHelper.FilterFunc != nil {
			opts = []cmp.Option{
				cmp.FilterPath(kindInfo.CompareHelper.FilterFunc, cmp.Comparer(kindInfo.CompareHelper.Comparer)),
			}
		}

		diff := cmp.Diff(res1.Obj, res2.Obj, opts...)
		fmt.Println(diff)
	}

}

func main() {

	files, err := os.ReadDir("/Users/afridshaikh/projects/golang/junk/prod/")
	if err != nil {
		panic("failed to read template")
	}

	for _, file := range files {

		prodObj, err := parser.ReadHelmTemplate(filepath.Join(PROD + file.Name()))
		if err != nil {
			panic("failed to parse prod. Error: " + err.Error())
		}

		prodDRObjs, err := parser.ReadHelmTemplate(filepath.Join(DR + file.Name()))
		if err != nil {
			panic("failed to parse prod-dr. Error: " + err.Error())
		}

		for key, prodRes := range prodObj {
			drRes, ok := prodDRObjs[key]
			if !ok {
				fmt.Printf("Resource %s missing in prod\n", key)
				continue
			}
			compareResources(prodRes, drRes)
		}
	}

}
