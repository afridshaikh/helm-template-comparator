package kindregistry

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	appsv1 "k8s.io/api/apps/v1"
)


func ImageIDPath(p cmp.Path) bool {
	return strings.Contains(p.String(), "Containers") &&
		strings.HasSuffix(p.String(), "Image")
}

func ImageComparer(a, b string) bool {
	normalize := func(img string) string {
		parts := strings.Split(img, ":")
		name := parts[0]
		if strings.HasSuffix(name, "-1") {
			name = strings.TrimSuffix(name, "-1")
		}
		if len(parts) > 1 {
			return name + ":" + parts[1]
		}
		return name
	}
	return normalize(a) == normalize(b)
}


func DeployCompareHelper() CompareHelper {
	return CompareHelper{
		FilterFunc: ImageIDPath,
		Comparer: ImageComparer,
	}
}

func DeployType() func () interface{} {
	return func () interface{} {return &appsv1.Deployment{}}
}


