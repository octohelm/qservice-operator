package deploy

import (
	qservice_operator "github.com/octohelm/qservice-operator/deploy/component"
)

qservice_operator & {
	#context: "hw-dev"

	#values: ingressGateways: {
		"auto-internal": "hw-dev.rktl.xyz"
		external:        "hw-dev.querycap.com"
	}
}
