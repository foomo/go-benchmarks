package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	// Mogno
	mlock "github.com/square/mongo-lock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// https://github.com/go-redsync/redsync
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	// github.com/bsm/redislock
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"

	// https://github.com/etcd-io/etcd
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func BenchmarkLockByMongoSquareLock(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(b, err)

	client := mlock.NewClient(db.Database("test").Collection(RandomString(5)))
	err = client.CreateIndexes(ctx)
	require.NoError(b, err)

	resourceName := RandomString(16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.XLock(ctx, resourceName, "test", mlock.LockDetails{})
		require.NoError(b, err)
		_, err = client.Unlock(ctx, "test")
		require.NoError(b, err)
	}

}

func BenchmarkLockByRedisRedSync(b *testing.B) {
	b.ReportAllocs()
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	mutex := rs.NewMutex("test")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := mutex.LockContext(ctx)
		require.NoError(b, err)
		_, err = mutex.UnlockContext(ctx)
		require.NoError(b, err)
	}
}

func BenchmarkLockByBSMRedisLock(b *testing.B) {
	b.ReportAllocs()

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})
	defer client.Close()

	// Create a new lock client.
	locker := redislock.New(client)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock, err := locker.Obtain(ctx, "my-key", 100*time.Millisecond, nil)
		require.NoError(b, err)
		err = lock.Release(ctx)
		require.NoError(b, err)
	}
}

func BenchmarkLockByETCD(b *testing.B) {
	b.ReportAllocs()

	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// create a sessions to aqcuire a lock
	s, _ := concurrency.NewSession(cli)
	defer s.Close()
	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := l.Lock(ctx)
		require.NoError(b, err)
		err = l.Unlock(ctx)
		require.NoError(b, err)
	}

}

func RandomString(len int) string {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}
