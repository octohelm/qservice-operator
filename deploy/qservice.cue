package deploy

import (
	"github.com/octohelm/cuem/release"
)

release.#Release & {
	#name:      "srv-test"
	#namespace: "default"
	#context:   "hw-sg"

	spec: {
		kube: {
			apiVersion: "serving.octohelm.tech/v1alpha1"
			kind:       "QService"
			metadata: name: #name
			metadata: annotations: {
				"autoscaling.octohelm.tech/maxScale": "5"
				"autoscaling.octohelm.tech/metrics":  "Resource(name = cpu, targetAverageUtilization = 70)"
				"autoscaling.octohelm.tech/minScale": "1"
			}
			spec: {
				envs: {
					K: "123"
				}
				image: "nginx:alpine"
				ports: ["80"]
				replicas: 2
			}
		}
	}
}
