# K8S Slackbot
A slack bot built to control kubernetes cluster.
based on https://github.com/danielqsj/k8s-slackbot 


Image
-------------
the build is multi staged
the first stage is based on go 1.18 image and builds the artifact

afterwards the next stage copies the artifact and all connecting dynamic libraries to a busybox image
total size of the final image: 11 MB

Arguments
-------------
- **kubecfg-file** (*string*): Location of kubecfg file for access to kubernetes master service;
- **bot-token** (*string*): Token of slack bot to use
- **debug** (*boolean*): Whether enable debug log

Usage
-------------
```
$ docker build -t k8s-slackbot .
$ docker run -v ~/.kube/config:/etc/kubernetes/kubeconfig k8s-slackbot --kubecfg-file=/etc/kubernetes/kubeconfig --bot-token=$(bot-token)
```
Then you can talk to your slack bot via slack direct message.
there are 2 commands:
- **list** - return a list of running pods, their age, and their version
- **logs** - takes 2 argument for the pod name and tail limit and return the number of tail logs from the pod 

**Enjoy it.**