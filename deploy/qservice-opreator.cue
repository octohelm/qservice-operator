package deploy

import (
	qservice_operator "github.com/octohelm/qservice-operator/deploy/component"
)

qservice_operator & {
	#context: "hw-sg"

	#values: ingressGateways: {
		"auto-internal": "hw-sg.rktl.xyz"
		external:        "hw-sg.querycap.com"
	}
}
