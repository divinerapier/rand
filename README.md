# Rand

A rand package copied from standard library with a thread safe type additional.

## Test

``` bash
# testing standard rand library isn't threads safe
$ go test -timeout 30s -run '^(TestStdConcurrent)$' -v

# testing locked rand is thread safe
$ go test -timeout 30s -run '^(TestLockedConcurrent)$' -v

# testing the behavior of locked rand is same to standard library
$ go test -timeout 30s -run '^(TestLockedSameAsStd)$' -v
```
