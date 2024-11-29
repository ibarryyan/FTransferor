## FTransferor 

#### Build

```shell
go mod tidy 

go build 
```

#### Usage

server：

```shell
./FTransferor server --path filepath --port 8081 
```
```markdown
- path: file save path, default tmp
- port: server port, default 8088
```

client：

```shell
./FTransferor.exe cli --server localhost:8088 --file .\tml\f.zip
```

```markdown
- server: server addr
- file: upload file name
```

#### TODO 

1. [x] support web
2. [x] support secret
