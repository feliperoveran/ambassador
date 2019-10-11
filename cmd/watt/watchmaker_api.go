package main

import (
	"fmt"
	"os"

	"github.com/datawire/ambassador/pkg/consulwatch"
	"github.com/datawire/ambassador/pkg/supervisor"
)

type WatchSet struct {
	KubernetesWatches []KubernetesWatchSpec         `json:"kubernetes-watches"`
	ConsulWatches     []consulwatch.ConsulWatchSpec `json:"consul-watches"`
}

// Interpolate values into specific watches in specific places. This is not a generic method but could be made one
// eventually if so desired by implementing interpolate() methods on the various types contained within the WatchSet
// struct.
//
// FIXES:
// 	- https://github.com/datawire/ambassador/issues/110
//	- https://github.com/datawire/ambassador/issues/1508
func (w *WatchSet) interpolate() WatchSet {
	result := WatchSet{KubernetesWatches: w.KubernetesWatches}

	if w.ConsulWatches != nil {
		modifiedConsulWatchSpecs := make([]consulwatch.ConsulWatchSpec, 0)
		for _, s := range w.ConsulWatches {
			modifiedConsulWatchSpecs = append(modifiedConsulWatchSpecs, consulwatch.ConsulWatchSpec{
				Id:            s.Id,
				ServiceName:   s.ServiceName,
				Datacenter:    s.Datacenter,
				ConsulAddress: os.ExpandEnv(s.ConsulAddress),
			})
		}

		result.ConsulWatches = modifiedConsulWatchSpecs
	}

	return result
}

type KubernetesWatchSpec struct {
	Kind          string `json:"kind"`
	Namespace     string `json:"namespace"`
	FieldSelector string `json:"field-selector"`
	LabelSelector string `json:"label-selector"`
}

func star(s string) string {
	if s == "" {
		return "*"
	} else {
		return s
	}
}

func (k KubernetesWatchSpec) WatchId() string {
	return fmt.Sprintf("%s|%s|%s|%s", k.Kind, star(k.Namespace), star(k.FieldSelector), star(k.LabelSelector))
}

// IKubernetesWatchMaker is an interface for KubernetesWatchMaker implementations. It mostly exists to facilitate the
// creation of testing mocks.
type IKubernetesWatchMaker interface {
	MakeKubernetesWatch(spec KubernetesWatchSpec) (*supervisor.Worker, error)
}

// IConsulWatchMaker is an interface for ConsulWatchMaker implementations. It mostly exists to facilitate the
// creation of testing mocks.
type IConsulWatchMaker interface {
	MakeConsulWatch(spec consulwatch.ConsulWatchSpec) (*supervisor.Worker, error)
}
