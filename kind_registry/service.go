package kindregistry

import corev1 "k8s.io/api/core/v1"


func ServiceCompareHelper() CompareHelper {
	return CompareHelper{
		FilterFunc: nil,
		Comparer: nil,
	}
}

func ServiceType() func () interface{} {
	return func () interface{} {return &corev1.Service{}}
}
