local k = import 'k/main.libsonnet';
local version = importstr '../.version';

{

  values:: {
    name: 'qservice-operator',
    namespace: 'qservice-operator',
    version: version,
    replicas: 1,
    imagePullPolicy: 'IfNotPresent',
    ingressGateways: 'auto-internal:hw-dev.rktl.xyz,external:hw-dev.querycap.com',
  },

  images:: {
    qserivce_operator: 'docker.io/octohelm/qservice-operator:' + version,
  },

  container::
    k.core.v1.container.new($.values.name, $.images.qserivce_operator) +
    k.core.v1.container.withImagePullPolicy($.values.imagePullPolicy) +
    k.core.v1.container.withArgs(['--enable-leader-election']) +
    k.core.v1.container.withEnv({
      name: 'INGRESS_GATEWAYS',
      value: $.values.ingressGateways,
    })
  ,

  deployment:
    k.apps.v1.deployment.new(
      name=$.values.name,
      replicas=$.values.replicas,
      containers=[$.container],
    ) +
    k.apps.v1.deployment.spec.template.metadata.withLabels({ app: $.values.name }) +
    k.apps.v1.deployment.spec.template.metadata.withAnnotations({ 'sidecar.istio.io/inject': 'false' }) +
    k.apps.v1.deployment.spec.selector.withMatchLabels({ app: $.values.name }) +
    k.apps.v1.deployment.spec.template.spec.withServiceAccount($.serviceAccount.metadata.name)
  ,

  rules:: [
    {
      verbs: ['*'],
      apiGroups: ['*'],
      resources: ['*'],
    },
    {
      verbs: ['*'],
      nonResourceURLs: ['*'],
    },
  ],

  serviceAccount:
    k.core.v1.serviceAccount.new($.values.name)
  ,

  clusterRole:
    k.rbac.v1.clusterRole.new($.values.name) +
    k.rbac.v1.clusterRole.withRules($.rules)
  ,

  clusterRoleBinding:
    k.rbac.v1.clusterRoleBinding.new($.values.name) +
    k.rbac.v1.clusterRoleBinding.withSubjects([{
      kind: $.serviceAccount.kind,
      name: $.serviceAccount.metadata.name,
      namespace: $.values.namespace,
    }]) +
    {
      roleRef: {
        apiGroup: 'rbac.authorization.k8s.io',
        kind: $.clusterRole.kind,
        name: $.clusterRole.metadata.name,
      },
    },
}
