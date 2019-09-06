# Rand

A rand package copied from standard library with a thread safe type additional.

## Test

``` bash
# test standard rand library isn't threads safe
$ go test -timeout 30s -run '^(TestStdConcurrent)$' -v

# test the behavior of locked rand is same to standard library
$ go test -timeout 30s -run '^(TestLockedSameAsStd)$' -v
```
