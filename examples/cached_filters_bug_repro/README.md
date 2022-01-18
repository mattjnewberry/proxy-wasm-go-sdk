## Cached Filters Bug Repro

### Expected Behaviour

The envoy proxy should fetch and run two different filters from the remoute source

### Current Behaviour

The envoy proxy fetches a single filter from the remote source and runs it twice

### Possible Solution

Not sure!

### Steps to Reprodue

```
$(cd ../../ &&  make build.example name=cached_filters_bug_repro)

# Verify the two binaries are distinct and update envoy.yaml with the correct sha256 values
openssl dgst -sha256 filter_one/main.go.wasm 
openssl dgst -sha256 filter_two/main.go.wasm

# Setup the remote source
docker build --no-cache -t filter-server .
docker run -d -p 80:8080 filter-server

# Verify the binaries served by the remote source are distinct
curl -s http://localhost/filter_one.go.wasm --output filter_one.go.wasm
curl -s http://localhost/filter_two.go.wasm --output filter_two.go.wasm
openssl dgst -sha256 filter_one.go.wasm
openssl dgst -sha256 filter_two.go.wasm

# Run the envoy proxy, you should see one of the filters logging x2
$(cd ../../ &&  make run name=helloworld)


```

### Logs

Debug startup logs of the envoy proxy

```
init manager Server initializing
init manager Server initializing target Listener-init-target main
init manager Listener-local-init-manager main 17126630817346794449 initializing
init manager Listener-local-init-manager main 17126630817346794449 initializing target RemoteAsyncDataProvider
fetch remote data from [uri = http://localhost/filter_one.go.wasm]: start
[C0][S12077228616084491646] cluster 'filter_server' match for URL '/filter_one.go.wasm'
[C0][S12077228616084491646] router decoding headers:
':path', '/filter_one.go.wasm'
':authority', 'localhost'
':method', 'GET'
':scheme', 'http'
'x-envoy-internal', 'true'
'x-forwarded-for', '172.20.164.123'
'x-envoy-expected-rq-timeout-ms', '60000'

queueing stream due to no available connections
trying to create new connection
creating a new connection
[C0] connecting
[C0] connecting to 127.0.0.1:80
[C0] connection in progress
init manager Listener-local-init-manager main 17126630817346794449 initializing target RemoteAsyncDataProvider
fetch remote data from [uri = http://localhost/filter_two.go.wasm]: start
[C0][S16741361052463840377] cluster 'filter_server' match for URL '/filter_two.go.wasm'
[C0][S16741361052463840377] router decoding headers:
':path', '/filter_two.go.wasm'
':authority', 'localhost'
':method', 'GET'
':scheme', 'http'
'x-envoy-internal', 'true'
'x-forwarded-for', '172.20.164.123'
'x-envoy-expected-rq-timeout-ms', '60000'

queueing stream due to no available connections
trying to create new connection
creating a new connection
[C1] connecting
[C1] connecting to 127.0.0.1:80
[C1] connection in progress
starting main dispatch loop
[C0] connected
[C0] connected
[C0] attaching to next stream
[C0] creating stream
[C0][S12077228616084491646] pool ready
[C1] connected
[C1] connected
[C1] attaching to next stream
[C1] creating stream
[C0][S16741361052463840377] pool ready
[C0][S12077228616084491646] upstream headers complete: end_stream=false
async http request response headers (end_stream=false):
':status', '200'
'server', 'nginx/1.21.4'
'date', 'Tue, 18 Jan 2022 14:16:26 GMT'
'content-type', 'application/wasm'
'content-length', '289434'
'last-modified', 'Tue, 18 Jan 2022 14:12:01 GMT'
'connection', 'keep-alive'
'etag', '"61e6cab1-46a9a"'
'accept-ranges', 'bytes'
'x-envoy-upstream-service-time', '5'

[C0][S16741361052463840377] upstream headers complete: end_stream=false
async http request response headers (end_stream=false):
':status', '200'
'server', 'nginx/1.21.4'
'date', 'Tue, 18 Jan 2022 14:16:26 GMT'
'content-type', 'application/wasm'
'content-length', '289267'
'last-modified', 'Tue, 18 Jan 2022 14:11:58 GMT'
'connection', 'keep-alive'
'etag', '"61e6caae-469f3"'
'accept-ranges', 'bytes'
'x-envoy-upstream-service-time', '6'

[C0] response complete
fetch remote data [uri = http://localhost/filter_one.go.wasm]: success
Base Wasm created 1 now active
Thread-Local Wasm created 2 now active
wasm log: Filter 1: OnPluginStart from Go!
~Wasm 1 remaining active
Thread-Local Wasm created 2 now active
wasm log: Filter 1: OnPluginStart from Go!
target RemoteAsyncDataProvider initialized, notifying init manager Listener-local-init-manager main 17126630817346794449
[C0] response complete
[C0] destroying stream: 0 remaining
[C1] response complete
fetch remote data [uri = http://localhost/filter_two.go.wasm]: success
target RemoteAsyncDataProvider initialized, notifying init manager Listener-local-init-manager main 17126630817346794449
init manager Listener-local-init-manager main 17126630817346794449 initialized, notifying Listener-local-init-watcher main
target Listener-init-target main initialized, notifying init manager Server
init manager Server initialized, notifying RunHelper
all dependencies initialized. starting workers
starting worker 0
starting worker 1
worker entering dispatch loop
worker entering dispatch loop
adding TLS cluster filter_server
completionThread running
completionThread running
membership update for TLS cluster filter_server added 1 removed 0
adding TLS cluster filter_server
membership update for TLS cluster filter_server added 1 removed 0
Thread-Local Wasm created 3 now active
Thread-Local Wasm created 4 now active
wasm log: Filter 1: OnPluginStart from Go!
wasm log: Filter 1: OnPluginStart from Go!
```