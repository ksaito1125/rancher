package node

import (
	"github.com/rancher/types/apis/management.cattle.io/v3"
	corev1 "k8s.io/api/core/v1"
)

func GetNodeName(node *v3.Node) string {
	if node.Status.NodeName != "" {
		return node.Status.NodeName
	}
	// to handle the case when node was provisioned first
	if node.Status.NodeConfig != nil {
		if node.Status.NodeConfig.HostnameOverride != "" {
			return node.Status.NodeConfig.HostnameOverride
		}
	}
	return ""
}

func IsNodeForNode(node *corev1.Node, machine *v3.Node) bool {
	nodeName := GetNodeName(machine)
	if nodeName == node.Name {
		return true
	}

	// search by rke external-ip annotations
	machineAddress := ""
	if machine.Status.NodeConfig != nil {
		if machine.Status.NodeConfig.InternalAddress == "" {
			// rke defaults internal address to address
			machineAddress = machine.Status.NodeConfig.Address
		} else {
			machineAddress = machine.Status.NodeConfig.InternalAddress
		}
	}

	if machineAddress == "" {
		return false
	}

	if machineAddress == getNodeInternalAddress(node) {
		return true
	}

	return false
}

func getNodeInternalAddress(node *corev1.Node) string {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return address.Address
		}
	}
	return ""
}