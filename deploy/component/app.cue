package component

import (
	"strings"

	"github.com/octohelm/cuem/release"
)

release.#Release & {
	#name:      "qservice-operator"
	#namespace: "\(#name)"

	spec: {
		serviceAccounts: "\(#name)": {
			#role: "ClusterRole"
			#rules: [
				{
					verbs: ["*"]
					apiGroups: ["*"]
					resources: ["*"]
				},
				{
					verbs: ["*"]
					nonResourceURLs: ["*"]
				},
			]
		}

		deployments: "\(#name)": {

			#containers: "qservice-operator": {
				image:           "\(#values.image.hub)/\(#values.image.name):\(#values.image.tag)"
				imagePullPolicy: "\(#values.image.pullPolicy)"
				args: ["--enable-leader-election"]

				#envVars: INGRESS_GATEWAYS: strings.Join([ for k, v in #values.ingressGateways {"\(k):\(v)"}], ",")
			}

			spec: template: spec: serviceAccountName: #name
		}
	}
}
