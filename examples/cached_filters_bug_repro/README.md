## Cached Filters Bug Repro

### Expected Behaviour

The envoy proxy should fetch and run two different filters from the remote source

### Current Behaviour

The envoy proxy fetches both envoy filters from the correct remote source (and ensures sha256 is correct) but loads the result of the initial fetch into envoy twice.

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
curl http://localhost/filter_one.go.wasm --output - | openssl dgst -sha256
curl http://localhost/filter_two.go.wasm --output - | openssl dgst -sha256

# Run the envoy proxy, you should see one of the filters logging x2
$(cd ../../ &&  make run name=cached_filters_bug_repro)


```

### Context

We're hosting out own binary server internally and found that when using remote source to load multiple envoy filters, we found instead of running both filters it would run a single filter twice.
It seems whichever remote data source envoy fetches first is then loaded twice.

We found a workaround where one of the filters is using a remote source and the other declares an extension config and uses config discovery. This implementation then runs both filters as expected.

### Logs

Debug startup logs of the envoy proxy

```
initializing epoch 0 (base id=0, hot restart version=disabled)
statically linked extensions:
  envoy.stats_sinks: envoy.dog_statsd, envoy.graphite_statsd, envoy.metrics_service, envoy.stat_sinks.dog_statsd, envoy.stat_sinks.graphite_statsd, envoy.stat_sinks.hystrix, envoy.stat_sinks.metrics_service, envoy.stat_sinks.statsd, envoy.stat_sinks.wasm, envoy.statsd
  envoy.resolvers: envoy.ip
  envoy.matching.action: composite-action, skip
  envoy.upstreams: envoy.filters.connection_pools.tcp.generic
  envoy.bootstrap: envoy.bootstrap.wasm, envoy.extensions.network.socket_interface.default_socket_interface
  envoy.quic.proof_source: envoy.quic.proof_source.filter_chain
  envoy.matching.common_inputs: envoy.matching.common_inputs.environment_variable
  envoy.health_checkers: envoy.health_checkers.redis
  envoy.retry_priorities: envoy.retry_priorities.previous_priorities
  envoy.dubbo_proxy.protocols: dubbo
  envoy.matching.http.input: request-headers, request-trailers, response-headers, response-trailers
  envoy.http.stateful_header_formatters: preserve_case
  envoy.filters.http: envoy.bandwidth_limit, envoy.buffer, envoy.cors, envoy.csrf, envoy.ext_authz, envoy.ext_proc, envoy.fault, envoy.filters.http.adaptive_concurrency, envoy.filters.http.admission_control, envoy.filters.http.alternate_protocols_cache, envoy.filters.http.aws_lambda, envoy.filters.http.aws_request_signing, envoy.filters.http.bandwidth_limit, envoy.filters.http.buffer, envoy.filters.http.cache, envoy.filters.http.cdn_loop, envoy.filters.http.composite, envoy.filters.http.compressor, envoy.filters.http.cors, envoy.filters.http.csrf, envoy.filters.http.decompressor, envoy.filters.http.dynamic_forward_proxy, envoy.filters.http.dynamo, envoy.filters.http.ext_authz, envoy.filters.http.ext_proc, envoy.filters.http.fault, envoy.filters.http.grpc_http1_bridge, envoy.filters.http.grpc_http1_reverse_bridge, envoy.filters.http.grpc_json_transcoder, envoy.filters.http.grpc_stats, envoy.filters.http.grpc_web, envoy.filters.http.header_to_metadata, envoy.filters.http.health_check, envoy.filters.http.ip_tagging, envoy.filters.http.jwt_authn, envoy.filters.http.local_ratelimit, envoy.filters.http.lua, envoy.filters.http.oauth2, envoy.filters.http.on_demand, envoy.filters.http.original_src, envoy.filters.http.ratelimit, envoy.filters.http.rbac, envoy.filters.http.router, envoy.filters.http.set_metadata, envoy.filters.http.squash, envoy.filters.http.tap, envoy.filters.http.wasm, envoy.grpc_http1_bridge, envoy.grpc_json_transcoder, envoy.grpc_web, envoy.health_check, envoy.http_dynamo_filter, envoy.ip_tagging, envoy.local_rate_limit, envoy.lua, envoy.rate_limit, envoy.router, envoy.squash, match-wrapper
  envoy.dubbo_proxy.serializers: dubbo.hessian2
  envoy.dubbo_proxy.filters: envoy.filters.dubbo.router
  envoy.compression.decompressor: envoy.compression.brotli.decompressor, envoy.compression.gzip.decompressor
  envoy.dubbo_proxy.route_matchers: default
  envoy.retry_host_predicates: envoy.retry_host_predicates.omit_canary_hosts, envoy.retry_host_predicates.omit_host_metadata, envoy.retry_host_predicates.previous_hosts
  envoy.internal_redirect_predicates: envoy.internal_redirect_predicates.allow_listed_routes, envoy.internal_redirect_predicates.previous_routes, envoy.internal_redirect_predicates.safe_cross_scheme
  envoy.http.cache: envoy.extensions.http.cache.simple
  envoy.wasm.runtime: envoy.wasm.runtime.null, envoy.wasm.runtime.v8
  envoy.quic.server.crypto_stream: envoy.quic.crypto_stream.server.quiche
  envoy.http.original_ip_detection: envoy.http.original_ip_detection.custom_header, envoy.http.original_ip_detection.xff
  envoy.thrift_proxy.transports: auto, framed, header, unframed
  envoy.compression.compressor: envoy.compression.brotli.compressor, envoy.compression.gzip.compressor
  envoy.tls.cert_validator: envoy.tls.cert_validator.default, envoy.tls.cert_validator.spiffe
  envoy.filters.listener: envoy.filters.listener.http_inspector, envoy.filters.listener.original_dst, envoy.filters.listener.original_src, envoy.filters.listener.proxy_protocol, envoy.filters.listener.tls_inspector, envoy.listener.http_inspector, envoy.listener.original_dst, envoy.listener.original_src, envoy.listener.proxy_protocol, envoy.listener.tls_inspector
  envoy.transport_sockets.downstream: envoy.transport_sockets.alts, envoy.transport_sockets.quic, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.starttls, envoy.transport_sockets.tap, envoy.transport_sockets.tls, raw_buffer, starttls, tls
  envoy.transport_sockets.upstream: envoy.transport_sockets.alts, envoy.transport_sockets.quic, envoy.transport_sockets.raw_buffer, envoy.transport_sockets.starttls, envoy.transport_sockets.tap, envoy.transport_sockets.tls, envoy.transport_sockets.upstream_proxy_protocol, raw_buffer, starttls, tls
  envoy.filters.udp_listener: envoy.filters.udp.dns_filter, envoy.filters.udp_listener.udp_proxy
  envoy.grpc_credentials: envoy.grpc_credentials.aws_iam, envoy.grpc_credentials.default, envoy.grpc_credentials.file_based_metadata
  envoy.clusters: envoy.cluster.eds, envoy.cluster.logical_dns, envoy.cluster.original_dst, envoy.cluster.static, envoy.cluster.strict_dns, envoy.clusters.aggregate, envoy.clusters.dynamic_forward_proxy, envoy.clusters.redis
  envoy.matching.input_matchers: envoy.matching.matchers.consistent_hashing, envoy.matching.matchers.ip
  envoy.tracers: envoy.dynamic.ot, envoy.lightstep, envoy.tracers.datadog, envoy.tracers.dynamic_ot, envoy.tracers.lightstep, envoy.tracers.opencensus, envoy.tracers.skywalking, envoy.tracers.xray, envoy.tracers.zipkin, envoy.zipkin
  envoy.resource_monitors: envoy.resource_monitors.fixed_heap, envoy.resource_monitors.injected_resource
  envoy.thrift_proxy.filters: envoy.filters.thrift.rate_limit, envoy.filters.thrift.router
  envoy.rate_limit_descriptors: envoy.rate_limit_descriptors.expr
  envoy.thrift_proxy.protocols: auto, binary, binary/non-strict, compact, twitter
  envoy.filters.network: envoy.client_ssl_auth, envoy.echo, envoy.ext_authz, envoy.filters.network.client_ssl_auth, envoy.filters.network.connection_limit, envoy.filters.network.direct_response, envoy.filters.network.dubbo_proxy, envoy.filters.network.echo, envoy.filters.network.ext_authz, envoy.filters.network.http_connection_manager, envoy.filters.network.kafka_broker, envoy.filters.network.local_ratelimit, envoy.filters.network.mongo_proxy, envoy.filters.network.mysql_proxy, envoy.filters.network.postgres_proxy, envoy.filters.network.ratelimit, envoy.filters.network.rbac, envoy.filters.network.redis_proxy, envoy.filters.network.rocketmq_proxy, envoy.filters.network.sni_cluster, envoy.filters.network.sni_dynamic_forward_proxy, envoy.filters.network.tcp_proxy, envoy.filters.network.thrift_proxy, envoy.filters.network.wasm, envoy.filters.network.zookeeper_proxy, envoy.http_connection_manager, envoy.mongo_proxy, envoy.ratelimit, envoy.redis_proxy, envoy.tcp_proxy
  envoy.formatter: envoy.formatter.req_without_query
  envoy.upstream_options: envoy.extensions.upstreams.http.v3.HttpProtocolOptions, envoy.upstreams.http.http_protocol_options
  envoy.access_loggers: envoy.access_loggers.file, envoy.access_loggers.http_grpc, envoy.access_loggers.open_telemetry, envoy.access_loggers.stderr, envoy.access_loggers.stdout, envoy.access_loggers.tcp_grpc, envoy.access_loggers.wasm, envoy.file_access_log, envoy.http_grpc_access_log, envoy.open_telemetry_access_log, envoy.stderr_access_log, envoy.stdout_access_log, envoy.tcp_grpc_access_log, envoy.wasm_access_log
  envoy.guarddog_actions: envoy.watchdog.abort_action, envoy.watchdog.profile_action
  envoy.request_id: envoy.request_id.uuid
HTTP header map info:
Unable to use runtime singleton for feature envoy.http.headermap.lazy_map_min_size
Unable to use runtime singleton for feature envoy.reloadable_features.header_map_correctly_coalesce_cookies
Unable to use runtime singleton for feature envoy.http.headermap.lazy_map_min_size
Unable to use runtime singleton for feature envoy.reloadable_features.header_map_correctly_coalesce_cookies
Unable to use runtime singleton for feature envoy.http.headermap.lazy_map_min_size
Unable to use runtime singleton for feature envoy.reloadable_features.header_map_correctly_coalesce_cookies
Unable to use runtime singleton for feature envoy.http.headermap.lazy_map_min_size
Unable to use runtime singleton for feature envoy.reloadable_features.header_map_correctly_coalesce_cookies
  request header map: 632 bytes: :authority,:method,:path,:protocol,:scheme,accept,accept-encoding,access-control-request-method,authentication,authorization,cache-control,cdn-loop,connection,content-encoding,content-length,content-type,expect,grpc-accept-encoding,grpc-timeout,if-match,if-modified-since,if-none-match,if-range,if-unmodified-since,keep-alive,origin,pragma,proxy-connection,referer,te,transfer-encoding,upgrade,user-agent,via,x-client-trace-id,x-envoy-attempt-count,x-envoy-decorator-operation,x-envoy-downstream-service-cluster,x-envoy-downstream-service-node,x-envoy-expected-rq-timeout-ms,x-envoy-external-address,x-envoy-force-trace,x-envoy-hedge-on-per-try-timeout,x-envoy-internal,x-envoy-ip-tags,x-envoy-max-retries,x-envoy-original-path,x-envoy-original-url,x-envoy-retriable-header-names,x-envoy-retriable-status-codes,x-envoy-retry-grpc-on,x-envoy-retry-on,x-envoy-upstream-alt-stat-name,x-envoy-upstream-rq-per-try-timeout-ms,x-envoy-upstream-rq-timeout-alt-response,x-envoy-upstream-rq-timeout-ms,x-forwarded-client-cert,x-forwarded-for,x-forwarded-proto,x-ot-span-context,x-request-id
  request trailer map: 136 bytes: 
  response header map: 432 bytes: :status,access-control-allow-credentials,access-control-allow-headers,access-control-allow-methods,access-control-allow-origin,access-control-expose-headers,access-control-max-age,age,cache-control,connection,content-encoding,content-length,content-type,date,etag,expires,grpc-message,grpc-status,keep-alive,last-modified,location,proxy-connection,server,transfer-encoding,upgrade,vary,via,x-envoy-attempt-count,x-envoy-decorator-operation,x-envoy-degraded,x-envoy-immediate-health-check-fail,x-envoy-ratelimited,x-envoy-upstream-canary,x-envoy-upstream-healthchecked-cluster,x-envoy-upstream-service-time,x-request-id
  response trailer map: 160 bytes: grpc-message,grpc-status
No overload action is configured for envoy.overload_actions.shrink_heap.
No overload action is configured for envoy.overload_actions.reduce_timeouts.
No overload action is configured for envoy.overload_actions.stop_accepting_connections.
No overload action is configured for envoy.overload_actions.reject_incoming_connections.
No overload action is configured for envoy.overload_actions.reduce_timeouts.
No overload action is configured for envoy.overload_actions.stop_accepting_connections.
No overload action is configured for envoy.overload_actions.reject_incoming_connections.
admin address: 0.0.0.0:8001
runtime: {}
loading tracing configuration
loading 0 static secret(s)
loading 1 cluster(s)
completionThread running
transport socket match, socket default selected for host with address 127.0.0.1:80
initializing Primary cluster filter_server completed
init manager Cluster filter_server contains no targets
init manager Cluster filter_server initialized, notifying ClusterImplBase
adding TLS cluster filter_server
membership update for TLS cluster filter_server added 1 removed 0
cm init: init complete: cluster=filter_server primary=0 secondary=0
maybe finish initialize state: 0
cm init: adding: cluster=filter_server primary=0 secondary=0
maybe finish initialize state: 1
maybe finish initialize primary init clusters empty: true
loading 1 listener(s)
listener #0:
begin add/update listener: name=main hash=14510100115403736852
use full listener update path for listener name=main hash=14510100115403736852
  filter #0:
    name: envoy.http_connection_manager
  config: {"@type":"type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager","stat_prefix":"ingress_http","http_filters":[{"typed_config":{"@type":"type.googleapis.com/udpa.type.v1.TypedStruct","type_url":"type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm","value":{"config":{"vm_config":{"code":{"remote":{"sha256":"53265bea87bd5a917b5a6f77b3819985b191fb645a8af60e8aa2f7c2355dd555","http_uri":{"uri":"http://localhost/filter_one.go.wasm","timeout":"60s","cluster":"filter_server"},"retry_policy":{"num_retries":3}}},"runtime":"envoy.wasm.runtime.v8"},"configuration":{"value":"Filter 1","@type":"type.googleapis.com/google.protobuf.StringValue"}}}},"name":"envoy.filters.http.wasm"},{"typed_config":{"@type":"type.googleapis.com/udpa.type.v1.TypedStruct","value":{"config":{"vm_config":{"code":{"remote":{"retry_policy":{"num_retries":3},"sha256":"09373f29acb1a58f53d8546195e5701e1df539a7c676194d5fdd1e5d98a08078","http_uri":{"uri":"http://localhost/filter_two.go.wasm","timeout":"60s","cluster":"filter_server"}}},"runtime":"envoy.wasm.runtime.v8"},"configuration":{"@type":"type.googleapis.com/google.protobuf.StringValue","value":"Filter 2"}}},"type_url":"type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm"},"name":"envoy.filters.http.wasm"},{"name":"envoy.filters.http.router"}],"route_config":{"virtual_hosts":[{"name":"local_service","routes":[{"direct_response":{"status":200,"body":{"inline_string":"example body\n"}},"match":{"prefix":"/"}}],"domains":["*"]}],"name":"local_route"},"codec_type":"AUTO"}
    http filter #0
added target RemoteAsyncDataProvider to init manager Listener-local-init-manager main 14510100115403736852
      name: envoy.filters.http.wasm
    config: {
 "@type": "type.googleapis.com/udpa.type.v1.TypedStruct",
 "type_url": "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm",
 "value": {
  "config": {
   "vm_config": {
    "code": {
     "remote": {
      "sha256": "53265bea87bd5a917b5a6f77b3819985b191fb645a8af60e8aa2f7c2355dd555",
      "http_uri": {
       "uri": "http://localhost/filter_one.go.wasm",
       "timeout": "60s",
       "cluster": "filter_server"
      },
      "retry_policy": {
       "num_retries": 3
      }
     }
    },
    "runtime": "envoy.wasm.runtime.v8"
   },
   "configuration": {
    "value": "Filter 1",
    "@type": "type.googleapis.com/google.protobuf.StringValue"
   }
  }
 }
}

    http filter #1
added target RemoteAsyncDataProvider to init manager Listener-local-init-manager main 14510100115403736852
      name: envoy.filters.http.wasm
    config: {
 "@type": "type.googleapis.com/udpa.type.v1.TypedStruct",
 "value": {
  "config": {
   "vm_config": {
    "code": {
     "remote": {
      "retry_policy": {
       "num_retries": 3
      },
      "sha256": "09373f29acb1a58f53d8546195e5701e1df539a7c676194d5fdd1e5d98a08078",
      "http_uri": {
       "uri": "http://localhost/filter_two.go.wasm",
       "timeout": "60s",
       "cluster": "filter_server"
      }
     }
    },
    "runtime": "envoy.wasm.runtime.v8"
   },
   "configuration": {
    "@type": "type.googleapis.com/google.protobuf.StringValue",
    "value": "Filter 2"
   }
  }
 },
 "type_url": "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm"
}

    http filter #2
      name: envoy.filters.http.router
    config: {}

new fc_contexts has 1 filter chains, including 1 newly built
added target Listener-init-target main to init manager Server
Create listen socket for listener main on address 0.0.0.0:18000
main: Setting socket options succeeded
Set listener main socket factory local address to 0.0.0.0:18000
add active listener: name=main, hash=14510100115403736852, address=0.0.0.0:18000
loading stats configuration
init manager RTDS contains no targets
init manager RTDS initialized, notifying RTDS
RTDS has finished initialization
continue initializing secondary clusters
maybe finish initialize state: 2
maybe finish initialize primary init clusters empty: true
maybe finish initialize secondary init clusters empty: true
maybe finish initialize cds api ready: false
cm init: all clusters initialized
there is no configured limit to the number of allowed active connections. Set a limit via the runtime key overload.global_downstream_max_connections
all clusters initialized. initializing init manager
init manager Server initializing
init manager Server initializing target Listener-init-target main
init manager Listener-local-init-manager main 14510100115403736852 initializing
init manager Listener-local-init-manager main 14510100115403736852 initializing target RemoteAsyncDataProvider
fetch remote data from [uri = http://localhost/filter_one.go.wasm]: start
[C0][S11246675459242952040] cluster 'filter_server' match for URL '/filter_one.go.wasm'
[C0][S11246675459242952040] router decoding headers:
':path', '/filter_one.go.wasm'
':authority', 'localhost'
':method', 'GET'
':scheme', 'http'
'x-envoy-internal', 'true'
'x-forwarded-for', '172.20.160.138'
'x-envoy-expected-rq-timeout-ms', '60000'

queueing stream due to no available connections
trying to create new connection
creating a new connection
[C0] connecting
[C0] connecting to 127.0.0.1:80
[C0] connection in progress
init manager Listener-local-init-manager main 14510100115403736852 initializing target RemoteAsyncDataProvider
fetch remote data from [uri = http://localhost/filter_two.go.wasm]: start
[C0][S2144931179243602233] cluster 'filter_server' match for URL '/filter_two.go.wasm'
[C0][S2144931179243602233] router decoding headers:
':path', '/filter_two.go.wasm'
':authority', 'localhost'
':method', 'GET'
':scheme', 'http'
'x-envoy-internal', 'true'
'x-forwarded-for', '172.20.160.138'
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
[C0][S11246675459242952040] pool ready
[C1] connected
[C1] connected
[C1] attaching to next stream
[C1] creating stream
[C0][S2144931179243602233] pool ready
[C0][S11246675459242952040] upstream headers complete: end_stream=false
async http request response headers (end_stream=false):
':status', '200'
'server', 'nginx/1.21.4'
'date', 'Wed, 02 Feb 2022 10:44:09 GMT'
'content-type', 'application/wasm'
'content-length', '292513'
'last-modified', 'Wed, 02 Feb 2022 10:36:09 GMT'
'connection', 'keep-alive'
'etag', '"61fa5e99-476a1"'
'accept-ranges', 'bytes'
'x-envoy-upstream-service-time', '4'

[C0][S2144931179243602233] upstream headers complete: end_stream=false
async http request response headers (end_stream=false):
':status', '200'
'server', 'nginx/1.21.4'
'date', 'Wed, 02 Feb 2022 10:44:09 GMT'
'content-type', 'application/wasm'
'content-length', '292513'
'last-modified', 'Wed, 02 Feb 2022 10:36:06 GMT'
'connection', 'keep-alive'
'etag', '"61fa5e96-476a1"'
'accept-ranges', 'bytes'
'x-envoy-upstream-service-time', '7'

[C0] response complete
fetch remote data [uri = http://localhost/filter_one.go.wasm]: success
Base Wasm created 1 now active
Thread-Local Wasm created 2 now active
wasm log: Filter 1 OnPluginStart: config: Filter 1, guid: 9457541284338866332, time: 1643798649875730000
~Wasm 1 remaining active
Thread-Local Wasm created 2 now active
wasm log: Filter 1 OnPluginStart: config: Filter 1, guid: 6567645268386784623, time: 1643798649880109000
target RemoteAsyncDataProvider initialized, notifying init manager Listener-local-init-manager main 14510100115403736852
[C0] response complete
[C0] destroying stream: 0 remaining
[C1] response complete
fetch remote data [uri = http://localhost/filter_two.go.wasm]: success
wasm log: Filter 1 OnPluginStart: config: Filter 2, guid: 5143901740323388853, time: 1643798649881247000
target RemoteAsyncDataProvider initialized, notifying init manager Listener-local-init-manager main 14510100115403736852
init manager Listener-local-init-manager main 14510100115403736852 initialized, notifying Listener-local-init-watcher main
target Listener-init-target main initialized, notifying init manager Server
init manager Server initialized, notifying RunHelper
all dependencies initialized. starting workers
starting worker 0
starting worker 1
worker entering dispatch loop
worker entering dispatch loop
adding TLS cluster filter_server
adding TLS cluster filter_server
completionThread running
completionThread running
membership update for TLS cluster filter_server added 1 removed 0
membership update for TLS cluster filter_server added 1 removed 0
Thread-Local Wasm created 3 now active
Thread-Local Wasm created 4 now active
wasm log: Filter 1 OnPluginStart: config: Filter 1, guid: 1788633596454771620, time: 1643798649886493000
wasm log: Filter 1 OnPluginStart: config: Filter 1, guid: 13500188827932050722, time: 1643798649886491000
wasm log: Filter 1 OnPluginStart: config: Filter 2, guid: 6513709612281615464, time: 1643798649886548000
wasm log: Filter 1 OnPluginStart: config: Filter 2, guid: 2787832971854728481, time: 1643798649886565000
[C1] response complete
[C1] destroying stream: 0 remaining
flushing stats
flushing stats
wasm log: Filter 1 OnTick: config: Filter 2, guid: 5143901740323388853, time: 1643798659843180000
wasm log: Filter 1 OnTick: config: Filter 1, guid: 6567645268386784623, time: 1643798659843212000
wasm log: Filter 1 OnTick: config: Filter 1, guid: 13500188827932050722, time: 1643798659886989000
wasm log: Filter 1 OnTick: config: Filter 2, guid: 2787832971854728481, time: 1643798659887021000
wasm log: Filter 1 OnTick: config: Filter 1, guid: 1788633596454771620, time: 1643798659890822000
wasm log: Filter 1 OnTick: config: Filter 2, guid: 6513709612281615464, time: 1643798659890843000
flushing stats
flushing stats
wasm log: Filter 1 OnTick: config: Filter 1, guid: 6567645268386784623, time: 1643798669843225000
wasm log: Filter 1 OnTick: config: Filter 2, guid: 5143901740323388853, time: 1643798669843257000
wasm log: Filter 1 OnTick: config: Filter 2, guid: 2787832971854728481, time: 1643798669890932000
wasm log: Filter 1 OnTick: config: Filter 1, guid: 13500188827932050722, time: 1643798669890964000
wasm log: Filter 1 OnTick: config: Filter 2, guid: 6513709612281615464, time: 1643798669892028000
wasm log: Filter 1 OnTick: config: Filter 1, guid: 1788633596454771620, time: 1643798669892049000
^Ccaught SIGINT
shutting down server instance
Notifying 0 callback(s) with completion.
main dispatch loop exited
worker exited dispatch loop
~Wasm 3 remaining active
Joining completionThread
completionThread exiting
Joined completionThread
shutting down thread local cluster manager
worker exited dispatch loop
~Wasm 2 remaining active
Joining completionThread
completionThread exiting
Joined completionThread
shutting down thread local cluster manager
flushing stats
ClusterImplBase destroyed
init manager Cluster filter_server destroyed
~Wasm 1 remaining active
~Wasm 0 remaining active
Joining completionThread
completionThread exiting
Joined completionThread
shutting down thread local cluster manager
[C1] closing data_to_write=0 type=1
[C1] closing socket: 1
[C1] disconnect. resetting 0 pending requests
[C1] client disconnected, failure reason: 
[C0] closing data_to_write=0 type=1
[C0] closing socket: 1
[C0] disconnect. resetting 0 pending requests
[C0] client disconnected, failure reason: 
exiting
RunHelper destroyed
destroying listener manager
destroying dispatcher worker_1
destroying dispatcher worker_0
Listener-local-init-watcher main destroyed
target RemoteAsyncDataProvider destroyed
target RemoteAsyncDataProvider destroyed
init manager Listener-local-init-manager main 14510100115403736852 destroyed
target Listener-init-target main destroyed
destroyed listener manager
destroying dispatcher workers_guarddog_thread
destroying dispatcher main_thread_guarddog_thread
destroying access logger /dev/null
destroyed access loggers
init manager RTDS destroyed
RTDS destroyed
destroying dispatcher main_thread
init manager Server destroyed

```