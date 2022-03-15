# Go Distributed Locking Benchmark


This is a test to benchmark existing go libraries for distributed locking for 

  - Redis: [github.com/go-redsync/redsync](https://github.com/go-redsync/redsync)
  - Redis: [github.com/bsm/redislock](https://github.com/bsm/redislock)
  - Mongo: [github.com/square/mongo-lock](https://github.com/square/mongo-lock)
  - ETCD:  [github.com/etcd-io/etcd](https://github.com/etcd-io/etcd)

## Run

To run the benchmarky you would need docker-compose and go.
After the dependencies are installed, run. 

``` make benchmark ```


## Results

The best performance was shown with redis, both redsync and redislock performing similarly.
Redislock has a better performance than redsync, but redsync has a nicer interface, reminding of a mutex lock.

```text
$ go test -bench=. -count=5 > out.txt

goos: darwin
goarch: amd64
pkg: distrocomp
cpu: VirtualApple @ 2.50GHz

BenchmarkLockByMongoSquareLock-10            452           2811117 ns/op           31242 B/op        525 allocs/op
BenchmarkLockByMongoSquareLock-10            492           2766211 ns/op           31338 B/op        525 allocs/op
BenchmarkLockByMongoSquareLock-10            415           2886690 ns/op           31318 B/op        525 allocs/op
BenchmarkLockByMongoSquareLock-10            421           2444429 ns/op           31279 B/op        525 allocs/op
BenchmarkLockByMongoSquareLock-10            375           2785962 ns/op           31264 B/op        525 allocs/op
BenchmarkLockByRedisRedSync-10              1080           1211415 ns/op            1616 B/op         39 allocs/op
BenchmarkLockByRedisRedSync-10              1054           1091482 ns/op            1621 B/op         39 allocs/op
BenchmarkLockByRedisRedSync-10              1219           1019288 ns/op            1619 B/op         39 allocs/op
BenchmarkLockByRedisRedSync-10              1041           1021278 ns/op            1621 B/op         39 allocs/op
BenchmarkLockByRedisRedSync-10              1311           1377340 ns/op            1616 B/op         39 allocs/op
BenchmarkLockByBSMRedisLock-10              1124           1264446 ns/op            1167 B/op         28 allocs/op
BenchmarkLockByBSMRedisLock-10               908           1152528 ns/op            1174 B/op         28 allocs/op
BenchmarkLockByBSMRedisLock-10              1077           1143109 ns/op            1162 B/op         28 allocs/op
BenchmarkLockByBSMRedisLock-10              1332            998186 ns/op            1171 B/op         28 allocs/op
BenchmarkLockByBSMRedisLock-10              1281            957423 ns/op            1168 B/op         28 allocs/op
BenchmarkLockByETCD-10                        18          71324093 ns/op           18776 B/op        309 allocs/op
BenchmarkLockByETCD-10                        19          69603401 ns/op           18658 B/op        307 allocs/op
BenchmarkLockByETCD-10                        19          52823399 ns/op           18381 B/op        303 allocs/op
BenchmarkLockByETCD-10                        21          73684962 ns/op           18592 B/op        307 allocs/op
BenchmarkLockByETCD-10                        19          55747954 ns/op           18963 B/op        313 allocs/op
```

```text
$ benchstat out.txt

name                      time/op
LockByMongoSquareLock-10  2.74ms ±11%
LockByRedisRedSync-10     1.14ms ±20%
LockByBSMRedisLock-10     1.10ms ±15%
LockByETCD-10             64.6ms ±18%

name                      alloc/op
LockByMongoSquareLock-10  31.3kB ± 0%
LockByRedisRedSync-10     1.62kB ± 0%
LockByBSMRedisLock-10     1.17kB ± 1%
LockByETCD-10             18.7kB ± 2%

name                      allocs/op
LockByMongoSquareLock-10     525 ± 0%
LockByRedisRedSync-10       39.0 ± 0%
LockByBSMRedisLock-10       28.0 ± 0%
LockByETCD-10                308 ± 2%
```



