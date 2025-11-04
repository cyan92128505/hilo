
### Generate ES key
```bash
// p-256
$ openssl ecparam -name prime256v1 -genkey -noout -out es256_private.pem
```
### Generate ES public key
```bash
// p-256 public
$ openssl ec -in es256_private.pem -pubout -out es256_public.pem
```

### wire injection
```go
wire.NewSet(NewEES256JWTFromOptions, wire.Bind(new(IJWT), new(*EES256JWT)))
```