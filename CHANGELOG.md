# V1.1.1 Fix judging file name

Before this release, if you gochk dirs including the file like .gotxt, you would be hung up.
This release fixes this problem.

# v1.1.0 Adding Exit Mode

If you want to exit with `1` when violations occur, please specify `-e=true` (default `false`):

```zsh
go run ../cmd/gochk/main.go -t=../../goilerplate -c=../configs/config.json -e=true
```

# v1.0.0 First Release
