package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/resource"
)

func DashboardManifest(folderUid string, dash dashboard.Dashboard) resource.Manifest {
	return resource.Manifest{
		ApiVersion: "dashboard.grafana.app/v1beta1",
		Kind:       "Dashboard",
		Metadata: resource.Metadata{
			Annotations: map[string]string{
				"grafana.app/folder": folderUid,
			},
			Name: *dash.Uid,
		},
		Spec: dash,
	}
}
