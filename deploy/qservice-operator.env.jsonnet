local q = import './qservice-operator.libsonnet';

{
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'demo',
  },
  spec: {
    namespace: q.values.namespace,
    apiServer: 'https://172.16.0.7:8443',
    injectLabels: true,
  },
  data: q,
}
