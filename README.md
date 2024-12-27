## FTransferor 

#### Run

```shell
/bin
```

#### Build

```shell
go mod tidy 

go build 
```

#### Usage

server：

```shell
./FTransferor server --path filepath --port 8081 --webport 8082 --secret D&J$HE23
```
```markdown
- path: file save path, default tmp
- port: server port, default 8088
- webport: web port, default 8089
- secret: secret key, default none
```

client：

```shell
./FTransferor.exe cli --server localhost:8088 --file .\tml\f.zip 

./FTransferor.exe cli --server localhost:8088 --action list --passwd D&J$HE23
```

```markdown
- server: server addr
- file: upload file name
- action: list or get
- passwd: secret key
```

#### TODO 

1. support sha256 checksum
