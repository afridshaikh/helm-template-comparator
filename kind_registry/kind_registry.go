package kindregistry

import "github.com/google/go-cmp/cmp"

type CompareHelper struct {
	FilterFunc func(cmp.Path) bool
	Comparer   func(a, b string) bool
}

type KindInfo struct {
	Type          func() interface{}
	CompareHelper CompareHelper
}


var KindRegistry = map[string]KindInfo{
	"Deployment": {Type: DeployType(), CompareHelper: DeployCompareHelper()},
	"Service": {Type: ServiceType(), CompareHelper: DeployCompareHelper()},
}
