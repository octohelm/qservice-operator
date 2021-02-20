{
  module: 'github.com/octohelm/qservice-operator',
  jpath: './jpath',
  replace: {
    k: 'github.com/jsonnet-libs/k8s-alpha/1.19',
  },
  require: {
    'github.com/jsonnet-libs/k8s-alpha':: 'v0.0.0-20210118111845-5e0d0738721f',
  },
}
