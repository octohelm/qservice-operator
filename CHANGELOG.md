# Change Log

All notable changes to this project will be documented in this file.
See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.



# [0.8.3](https://github.com/octohelm/qservice-operator/compare/v0.8.2...v0.8.3)

### Bug Fixes

* **fix** auto label role for deployment by service name ([15d4dfb](https://github.com/octohelm/qservice-operator/commit/15d4dfb5dbd789990bbc3d2f8dd245e302962686))



# [0.8.2](https://github.com/octohelm/qservice-operator/compare/v0.8.1...v0.8.2)



# [0.8.1](https://github.com/octohelm/qservice-operator/compare/v0.8.0...v0.8.1)

### Bug Fixes

* **fix** crd plural -es ([68385db](https://github.com/octohelm/qservice-operator/commit/68385db5c3594a1329f6de2828479debcb492beb))



# [0.8.0](https://github.com/octohelm/qservice-operator/compare/v0.7.3...v0.8.0)

### Features

* **feat** QIngress and QEgress ([f51b17e](https://github.com/octohelm/qservice-operator/commit/f51b17e3e5222a545d518aabbf857d1945eac0fc))



# [0.7.3](https://github.com/octohelm/qservice-operator/compare/v0.7.2...v0.7.3)

### Bug Fixes

* **fix** crd ([fb7540e](https://github.com/octohelm/qservice-operator/commit/fb7540ea89ab2a744d6fcba1f17a179140ea4d51))



# [0.7.2](https://github.com/octohelm/qservice-operator/compare/v0.7.1...v0.7.2)

### Bug Fixes

* **fix** crd need schema ([01e12df](https://github.com/octohelm/qservice-operator/commit/01e12df8eba1406ed5d87ab4496e06967cababd7))



# [0.7.1](https://github.com/octohelm/qservice-operator/compare/v0.7.0...v0.7.1)

### Bug Fixes

* **fix** kubectl.kubernetes.io/restartedAt changed should patch ([671cfef](https://github.com/octohelm/qservice-operator/commit/671cfef2ec952bd633e028ab05e634f26d317ad9))
* **fix** label gateway ([237c536](https://github.com/octohelm/qservice-operator/commit/237c536839b54dd957533dc944524043d8aae605))



# [0.7.0](https://github.com/octohelm/qservice-operator/compare/v0.6.1...v0.7.0)

### Bug Fixes

* **fix** ingress to virtual service support path type exact ([13c8ea9](https://github.com/octohelm/qservice-operator/commit/13c8ea92b7f46f162bfcc17ec85a17d789f45b58))


### Features

* **feat** abstract ingress gateways ([d6d0755](https://github.com/octohelm/qservice-operator/commit/d6d075507e03d51b03f86c81532258fe92f82bff))



# [0.6.1](https://github.com/octohelm/qservice-operator/compare/v0.6.0...v0.6.1)

### Bug Fixes

* **fix** crd should apply by k8s client directly ([07c9c52](https://github.com/octohelm/qservice-operator/commit/07c9c524ee40b1363e31319a75ad277b2f4d4823))



# [0.6.0](https://github.com/octohelm/qservice-operator/compare/v0.5.4...v0.6.0)

### Features

* **feat** crds automately apply ([8339cb4](https://github.com/octohelm/qservice-operator/commit/8339cb4441bdde73939bc30bd5cc83c93a74a0f7))



# [0.5.4](https://github.com/octohelm/qservice-operator/compare/v0.5.3...v0.5.4)

### Bug Fixes

* **fix** default labels of ingress ([8dd2410](https://github.com/octohelm/qservice-operator/commit/8dd2410cadd6ee7fcdedf00cf9dd00909bfa0cb1))



# [0.5.3](https://github.com/octohelm/qservice-operator/compare/v0.5.2...v0.5.3)



# [0.5.2](https://github.com/octohelm/qservice-operator/compare/v0.5.1...v0.5.2)

### Bug Fixes

* **fix** sync related ingress as QService status ([60216b9](https://github.com/octohelm/qservice-operator/commit/60216b9dffe47abbbb28cbdfb6adedc3fcfc0188))



# [0.5.1](https://github.com/octohelm/qservice-operator/compare/v0.5.0...v0.5.1)



# [0.5.0](https://github.com/octohelm/qservice-operator/compare/v0.4.0...v0.5.0)

### Bug Fixes

* **fix** crd must contains subresources with status for sync status ([70bdf79](https://github.com/octohelm/qservice-operator/commit/70bdf79cf2c2b5f7e1d98c29ce2b794db89707f0))


### Features

* **feat** drop ingresses and auto ingress from service with ClusterIP ([c72cb23](https://github.com/octohelm/qservice-operator/commit/c72cb23396fe2f2c1bc0ce0038a1e9768a794e2b))



# [0.4.0](https://github.com/octohelm/qservice-operator/compare/v0.3.1...v0.4.0)

### Features

* **feat** refactor to own by directly related ([b8c8d87](https://github.com/octohelm/qservice-operator/commit/b8c8d87e6d0868a405c0b2706bd61496c7d9e8fb))



# [0.3.1](https://github.com/octohelm/qservice-operator/compare/v0.3.0...v0.3.1)



# [0.3.0](https://github.com/octohelm/qservice-operator/compare/v0.2.3...v0.3.0)

### Features

* **feat** value format for value ref in env var ([5bdad1c](https://github.com/octohelm/qservice-operator/commit/5bdad1c17d22c8c994b4307842b36ec19cd8d139))



# [0.2.3](https://github.com/octohelm/qservice-operator/compare/v0.2.2...v0.2.3)

### Bug Fixes

* **fix** labels lens limit fix ([a004813](https://github.com/octohelm/qservice-operator/commit/a004813b89c1c8dcb723159e91a7995ae17679c0))



# [0.2.2](https://github.com/octohelm/qservice-operator/compare/v0.2.1...v0.2.2)

### Bug Fixes

* **fix** name of vs should be shorted ([f9e44da](https://github.com/octohelm/qservice-operator/commit/f9e44dad90f213db19405a2d83713978d527bdbb))



# [0.2.1](https://github.com/octohelm/qservice-operator/compare/v0.2.0...v0.2.1)

### Bug Fixes

* **fix** support multi auto ingress hosts ([65289d4](https://github.com/octohelm/qservice-operator/commit/65289d4cbe00da060707c88dd112d631b7bfb8f9))



# [0.2.0](https://github.com/octohelm/qservice-operator/compare/v0.1.2...v0.2.0)

### Features

* **feat** mv auto ingress host option to namespace to cluster ([c6f9191](https://github.com/octohelm/qservice-operator/commit/c6f91917d2d1f2e7ac4d7bc2dd7644364bf86085))



# [0.1.2](https://github.com/octohelm/qservice-operator/compare/v0.1.1...v0.1.2)

### Bug Fixes

* **fix** image pull secret fix ([f52a5b0](https://github.com/octohelm/qservice-operator/commit/f52a5b06df1dd2bbf88e25675295a8ea01417481))



# [0.1.1](https://github.com/octohelm/qservice-operator/compare/v0.1.0...v0.1.1)

### Bug Fixes

* **fix(controller/deployment):** only work in namespace with label autoscaling=enabled ([cce11b3](https://github.com/octohelm/qservice-operator/commit/cce11b35974ea803643b7a4b32c096482798aac8))



# [0.1.0](https://github.com/octohelm/qservice-operator/compare/v0.0.0...v0.1.0)

### Features

* **feat(controller/deployment):** to handle annotaions autoscaling.octohelm.tech/* to setup hpa ([6f2b740](https://github.com/octohelm/qservice-operator/commit/6f2b7400afea62bc848e79658d57652269410845))



# 0.0.0

### Features

* **feat** added Strategy ([69cf6d8](https://github.com/octohelm/qservice-operator/commit/69cf6d8f6ae68ae8786eebb6ae12de5de3b5bf0f))
