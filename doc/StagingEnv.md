#How to use staging environment for integration testing

## Staging environment overview

We design the staging environment to redirect the request specific traffic from the certain client to specific backend, so we can do integration testing in real production alike environment(using the same databse of production environment, share the same rest api url), no impact to the rest of world.  

Staging environment was composed by 2 parts

### Frontend pod

Because k8s now can't forward client real ip if we use k8s proxy, so we have to put our frontend program which is ip sensitive to a specific pod. In staging environment, the frontend pod's name called `ds2-staging-pod`, this pod handle the following services: 
```
- inappdata
- sharelinkfront
- urlgenerator
- dsaction
- counter
```

### Backend bundle

We seperated the backend into two parts, `deepshare2-match`(services included: match) and `deepshare2`(services included: cookie, appinfo).

## How to trace logs

In the new k8s cluster, we will no longer need to login to workstation and node of k8s cluster to trace the log, we can print the log just from your laptop.

### Step1: install docker-machine on your Mac

Before we begin to print log of your program, we need get `docker-machine` installed on your Mac, [click here](https://docs.docker.com/mac/step_one/) to see how to do it.

### Step 2: Login your local VM

Once `docker-machine` was installed on your Mac, a default vritual machine will be defined(`default`), now you can start it by `docker-machine start default` and login the VM using `docker-machine ssh default`.

### Step 4: Get CA file from admin

In the new k8s cluster, kubectl will use CA files to connect with k8s apiserver, so you need get the CA files, put them in sub directory `default-CA`.

### Step 5:  Get ready to use `kubectl`, and trace your log here

Now, you can just run command `docker run --rm -it -v ${PWD}:/mnt/qd r.fds.so:5000/qdkubectl:v1beta2` to startup and entering a container which you can use `kubectl` here.

Staging env was put into namespace `Staging`, so you should switch to staging by `kubectl config use-context staging`.

In staging namespace, you can list all running pods by `kubectl get po`.

You can print logs of pods by `kubectl logs <podname>`.

Running command `kubectl attach <podname>` will let you trace the logs interactively, just like '-f' option we used before. If you want to abort from the pod, you can just press ctrl-c.

## TODO
