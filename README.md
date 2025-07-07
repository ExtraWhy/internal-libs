# Shared internal libs
1. Config using yaml
2. Database module with sqlite3(experimental)
3. Player entity shared and enriched between service
4. Zapper logger based on uber [zap log](https://betterstack.com/community/guides/logging/go/zap/) 

# protobuf model (proto-models is just an example)
1. design your model like `players.proto`
2. go to proto-models
3. run `make.sh`
4. commit to git generated files form `players.proto`
5. Backend's `make.sh` shall do the magic for you

# DynConfig
If you'd like to regenerate the mocks, you'll need to run through the following steps:

1. install gomock/mockgen:

```
go install go.uber.org/mock/mockgen@v0.5.2
```

Then put the following line in your `.bashrc` file:

```
export PATH=$PATH:$(go env GOPATH)/bin
```

2. Run
```
make gogenerate
```

In order to run the unit tests:
```
make test
```

