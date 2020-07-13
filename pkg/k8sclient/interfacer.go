package k8sclient

import (
	v1 "k8s.io/api/core/v1"
)

type Interfacer interface {
	ListNode(Labelselector string) (*v1.NodeList, error)

	ListPod() (*v1.PodList, error)
}
