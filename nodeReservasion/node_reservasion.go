package nodeReservasion

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"time"
)

// Name of the plugin used in the plugin registry and configurations.
const Name = "NodeReservasion"

// NodeReservasion block pod bind until container checkpoint succeed.
type NodeReservasion struct {
	handle framework.FrameworkHandle
}

var _ framework.PermitPlugin = &NodeReservasion{}

// New creates a NodeReservasion.
func New(_ *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	return &NodeReservasion{handle: handle}, nil
}

// Name returns the name of the plugin.
func (nr *NodeReservasion) Name() string {
	return Name
}

func (nr *NodeReservasion) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {
	klog.Infoln("enter permit phase")
	//if pod with annotation["migrate-block"]=true, block it until annotation["migrate-block"]=false
	nr.handle.IterateOverWaitingPods(func(waitingPod framework.WaitingPod) {
		if block, ok := waitingPod.GetPod().Annotations["migrate-block"]; ok {
			if block == "false" {
				klog.Infoln("p.Annotations[\"migrate-block\"] = false, so allow the waiting pod")
				waitingPod.Allow(nr.Name())
			}
		}
	})

	if block, ok := p.Annotations["migrate-block"]; ok {
		if block == "true" {
			klog.Infoln("p.Annotations[\"migrate-block\"] = true, so block the creation of pod")
			return framework.NewStatus(framework.Wait, "Wait checkpoint files transporting complete!"), (15 * time.Minute)
		}
	}
	return nil, 0
}
