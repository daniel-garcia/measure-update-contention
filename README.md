```bash
dgarcia@host measure-update-contention % go run  main.go -concurrency 1 -iterations 100000
2024/10/10 14:58:50 Table foo recreated and initialized with one row.

Latency Distribution (microseconds):
Min: 124 µs
Max: 39423 µs
Mean: 461.34678 µs
P50: 253 µs
P90: 463 µs
P99: 5527 µs
dgarcia@host measure-update-contention % go run  main.go -concurrency 1 -iterations 100000
2024/10/10 15:00:29 Table foo recreated and initialized with one row.

Latency Distribution (microseconds):
Min: 124 µs
Max: 38271 µs
Mean: 452.59015 µs
P50: 253 µs
P90: 446 µs
P99: 4987 µs
dgarcia@host measure-update-contention % go run  main.go -concurrency 10 -iterations 10000
2024/10/10 15:02:05 Table foo recreated and initialized with one row.

Latency Distribution (microseconds):
Min: 170 µs
Max: 339967 µs
Mean: 7014.9937 µs
P50: 1853 µs
P90: 22127 µs
P99: 61759 µs
```

