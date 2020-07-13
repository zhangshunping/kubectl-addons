package utils

var AnnotationMap = map[string]map[string]string{
	"CA-NoDown": {"cluster-autoscaler.kubernetes.io/scale-down-disabled":"true"}, // cluster-autoscaler组件 根据着annotation决定 缩不缩减
}



