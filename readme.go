

找到对应node
显示详细信息

kubectl run --generator=run-pod/v1 tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash
[sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin)
[docker-debug](https://github.com/zeromake/docker-debug)
[docker api](https://godoc.org/github.com/docker/docker/client)
[docker api](https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190920-pod-pid-namespace.md)
[docker api](https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190920-pod-pid-namespace.md)
[Docker 核心技术与实现原理](https://draveness.me/docker)

`kubectl-debug` is an out-of-tree solution for [troubleshooting running pods](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/troubleshoot-running-pods.md)
https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190212-ephemeral-containers.md


// 核心的命令
kubectl attach POD -c CONTAINER


kubectl attach pod echoserver-64d97db464-8mrx8 -c nicolaka/netshoot

kubectl run --generator=run-pod/v1 tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash

joined 容器

joined 容器是另一种实现容器间通信的方式。它可以使两个或多个容器共享一个网络栈，共享网卡和配置信息，joined 容器之间可以通过 127.0.0.1 直接通信。

https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/20190920-pod-pid-namespace.md

ShareProcessNamespace