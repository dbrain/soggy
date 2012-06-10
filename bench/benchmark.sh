#!/bin/bash
# This may die after 16000~ requests due to ephermeral port sadness. On Mac OSX:
sudo sysctl -w net.inet.tcp.msl=1000
HOST="127.0.0.1"
WARMUP_HITS="10"
WARMUP_CONCURRENCY="1"
HITS="1000000"
CONCURRENCY="100"

echo "Warming up"
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9994/ > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/handler > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/html/fun/stuff > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/htmlNoParams/fun/stuff > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/json/fun/stuff > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/template/fun/stuff > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/self > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9991/bundle/admin > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9992/ > /dev/null
ab -n $WARMUP_HITS -c $WARMUP_CONCURRENCY http://$HOST:9993/ > /dev/null

echo "Basic net/http Server For Comparisons"
ab -n $HITS -c $CONCURRENCY http://$HOST:9994/ > results/basicserver.bench

echo "Single - Handler"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/handler > results/single-handler.bench

echo "Single - HTML response"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/html/fun/stuff > results/html-response.bench

echo "Single - HTML response - No params"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/htmlNoParams/fun/stuff > results/htmlNoParams-response.bench

echo "Single - JSON response"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/json/fun/stuff > results/json-response.bench

echo "Single - Template response"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/template/fun/stuff > results/template-response.bench

echo "Single - Self response handling"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/self > results/self-response.bench

echo "Single - Bundle with admin validation"
ab -n $HITS -c $CONCURRENCY http://$HOST:9991/bundle/admin > results/admin-response.bench

echo "Multiple Mounts - Last server"
ab -n $HITS -c $CONCURRENCY http://$HOST:9992/ > results/lastserver-response.bench

echo "Multiple Middleware"
ab -n $HITS -c $CONCURRENCY http://$HOST:9993/ > results/middleware-response.bench
