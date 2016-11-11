FROM r.fds.so:5000/golang1.5.3

RUN apt-get update && apt-get install -y --no-install-recommends vim
ADD . /go/src/github.com/MISingularity/deepshare2
WORKDIR /go/src/github.com/MISingularity/deepshare2
RUN git log -1 > /go/bin/deepshare2_version.txt

# keep files that will be read in runtime in /tmp dir
RUN cp -r /go/src/github.com/MISingularity/deepshare2/frontend/sharelink/js_template /tmp/js_template
RUN mkdir -p /tmp/ua_info
RUN cp -r /go/src/github.com/MISingularity/deepshare2/deepshared/uainfo/*.yaml /tmp/ua_info/

# build binary & remove source code
RUN godep go build -o /go/bin/deepshare2d /go/src/github.com/MISingularity/deepshare2/cmd/deepshared/main.go
RUN rm -rf /go/src

WORKDIR /go/bin

# restore files from /tmp dir
RUN mkdir -p /go/src/github.com/MISingularity/deepshare2/frontend/sharelink
RUN mkdir -p /tmp/ua_info /go/src/github.com/MISingularity/deepshare2/deepshared
RUN mv /tmp/js_template /go/src/github.com/MISingularity/deepshare2/frontend/sharelink/js_template
RUN mv /tmp/ua_info /go/src/github.com/MISingularity/deepshare2/deepshared/uainfo

CMD /go/bin/deepshare2d

## full list of arguments:
#   *argument name*     *default value* (*description*)
#   service-types       appcookiedevice,devicecookie,match,counter,dsaction,appinfo,token,sharelinkfront,urlgenerator,jsapi,inappdata,binddevicetocookie,dsusage
#   http-listen         0.0.0.0:8080
#   http-key            null    (http-key file name)
#   http-cert           null    (http-cert file name)
#   worker-id           0
#   data-center-id      0
#   redis-url           127.0.0.1:6379  (redis-url is required for dsusage. for other services, it is ignored when redis-sentinel-urls or redis-cluster-node-url is set)
#   redis-sentinel-urls 127.0.0.1:26379 (if there are multiple sentinels: "url1,url2,url3", it is ignored when redis-cluster-node-url is set)
#   redis-master-name   mymaster
#   redis-cluster-node-url  null (One of the redis 3.0 cluster nodes url)
#   redis-password      null
#   redis-pool-size     100             (size of redis connection pool)
#   nsq-url             127.0.0.1:4150  (nsqd tcp port)
#   match-url           127.0.0.1:8080
#   cookie-url          127.0.0.1:8080
#   appcookie-url       127.0.0.1:8080
#   appinfo-url         127.0.0.1:8080
#   urlgenerator-url    127.0.0.1:8080
#   genurl-base         http://127.0.0.1:8080
#   log-level           error           (log level to print: debug, info, error)

# Different services need different sets of arguments, listed as follows:
#   appcookiedevice:    http-listen,log-level,nsq-url,<redis-config>
#   devicecookie:       http-listen,log-level,nsq-url,<redis-config>
#   match:              http-listen,log-level,nsq-url,<redis-config>
#   counter:            http-listen,log-level,nsq-url
#   dsaction:           http-listen,log-level,nsq-url
#   appinfo:            http-listen,log-level,nsq-url,<redis-config>
#   token:              http-listen,log-level,nsq-url,data-center-id,worker-id
#   urlgenerator:       http-listen,log-level,nsq-url,<redis-config>,genurl-base,token-url
#   sharelinkfront:     http-listen,log-level,nsq-url,<redis-config>,genurl-base,match-url,appcookie-url,appinfo-url,token-url
#   jsapi:              http-listen,log-level,nsq-url,<redis-config>,genurl-base,match-url,appcookie-url,appinfo-url,token-url
#   inappdata:          http-listen,log-level,nsq-url,urlgenerator-url,match-url,cookie-url,appcookie-url
#   binddevicetocookie: http-listen,log-level,nsq-url,cookie-url,token-url
#   dsusage:            http-listen,log-level,<redis-config>

# <redis-config> can be one of the following sets of arguments
# 1) redis-url, redis-pool-size
#        to use single-instance redis
# 2) redis-sentinel-urls, redis-master-name, redis-pool-size
#        to use redis sentinel
# 3) redis-cluster-node-url, redis-pool-size
#        to use redis 3.0 cluster

EXPOSE 8080
