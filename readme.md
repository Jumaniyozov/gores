### Step 1. Initialize go mod
```
go mod init github.com/vlasove/go2/5.StandardWebServer
```

### Step 2. Where to find standard patterns?
**Useful link**: https://github.com/golang-standards/project-layout (There can be found information about structuring/packaging/refactoring of any Go apps)

### Step 3. Create an entry point for app
Standard pattern of entry point :
```
cmd/<app_name>/main.go
```
Here was created :
```
cmd/api/main.go
```

### Step 3. Initialize core of server
Standard pattern dictated in a following way
```
internal/app/<app_name>/<app_name>.go
```
We have ```internal/app/api/api.go```

### Step 4. Important point about configuration
**Rule**: in go:
* configurations are always stored in external files (.toml, .env)
* in Go projects always exists default configurations (exclusion - DB is intended not to have defaults)

### Step 5. Configuration of API server
Basically, for configuration only PORT is needed.
```
intrenal/app/api/config.go
```

### Step 6. Create configs
```
configs/<app_name>.toml or configs/.env
```

```
//api.toml
bind_addr = ":8080"
```

### Step 7. How to pass configurations?
We would want to pass the following way:
```
api.exe -path configs/api.toml
```

### Step 8. Configuration of http server
```
go get -u github.com/gorilla/mux
```
