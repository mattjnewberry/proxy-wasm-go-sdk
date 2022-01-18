FROM nginx:1.21.4

COPY default.conf /etc/nginx/conf.d/default.conf

ADD filter_one/main.go.wasm /usr/share/nginx/wasm/filter_one.go.wasm
ADD filter_two/main.go.wasm /usr/share/nginx/wasm/filter_two.go.wasm