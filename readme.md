// 调整代码

// kubectl run --generator=run-pod/v1 tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash
// [sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin)

// `kubectl-debug` is an out-of-tree solution for [troubleshooting running pods](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/troubleshoot-running-pods.md)
// https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190212-ephemeral-containers.md

// // 核心的命令
// kubectl attach POD -c CONTAINER
// kubectl attach pod echoserver-64d97db464-8mrx8 -c nicolaka/netshoot

// kubectl run --generator=run-pod/v1 tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash

// joined 容器

// https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190920-pod-pid-namespace.md

// ShareProcessNamespace