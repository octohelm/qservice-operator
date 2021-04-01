package component

#values: {
	ingressGateways: {[string]: string} & {
		"auto-internal": "hw-dev.rktl.xyz"
	}

	image: {
		hub:        *"docker.io/octohelm" | string
		name:       *"qservice-operator" | string
		tag:        *"0.10.1" | string
		pullPolicy: *"IfNotPresent" | string
	}
}
