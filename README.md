# sqlettuce

Redis server with sqlite backed support.

## Benchmarks

```
# run the benchmarks for the client
go test -bench . -run ^$
# run the benchmarks for all support `redis-benchmark` commands
redis-benchmark -p 6379 -q -c 10 -r 10000 -t "$(awk -F '"' '/case / { gsub(/, /, "\n"); print $2 }' commands/handle.go | tr ',' '\n' | sed '/^$/d' | tr '\n' ',' | sed 's/,$/\n/')"
```
