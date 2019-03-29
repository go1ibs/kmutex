# kmutex - lock by key mutex

[![Build Status](https://travis-ci.org/go1ibs/kmutex.svg)](https://travis-ci.org/go1ibs/kmutex)

```
ctx := context.Background()
kmux := New()

go func() {
	key := "green"
	kmux.Lock(ctx, key)
		...
	kmux.Unlock(ctx, key)
}()

go func() {
	key := "red"
	kmux.Lock(ctx, key)
		...
	kmux.Unlock(ctx, key)
}()

```
