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



## Database connection and migration schemes

### Step 9. Libraries to work with databases
```database/sql```
```sqlx```
```gosql```

### Step 10. Initialize database
```storage/storage.go```
Purpose of this model is:
* Instance of DB
* constructor of DB
* public method Open (setup connection)
* public method Close (close connection)


### Step 11. Initialize Storage
```storage.go```
The main problem lies inside the Open method, because in fact the low-level sql.Open is "lazy" (establishes a connection to the database only when the first query is made)

```config.go```
Contains a config instance and a constructor. The config attribute is only a connection string of the form :
```
"host=localhost port=5432 user=postgres password=postgres dbname=restapi sslmode=disable"
```

### Step 12. Add DB to API
Add new attribute storage
```
//Base API server instance description
type API struct {
	//UNEXPORTED FIELD!
	config *Config
	logger *logrus.Logger
	router *mux.Router
	storage *storage.Storage
}
```

Add new configurator:
```
//Configure storage (storage API)
func (a *API) configreStorageField() error {
	storage := storage.New(a.config.Storage)
	if err := storage.Open(); err != nil {
		return err
	}
	a.storage = storage
	return nil
}

```

### Step 13. Initial migration

#### For windows:
First install Scoop ```scoop```
* Open PowerShell: ```Set-ExecutionPolicy RemoteSigned -scope CurrentUser``` and ```Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://get.scoop.sh')```

After installation ```scoop``` run: ```scoop install migrate```

#### For linux 
* Run ```$ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz```
* Then move it to GOPATH ```mv migrate.linux-amd64 $GOPATH/bin/migrate```

### Step 13.1 Create migration repository
This repository will hold up/down pairs of sql migration requests to the database.
```
migrate create -ext sql -dir migrations UsersCreationMigration
```

### Step 13.2 Create up/down sql files
Look ```migrations/....up.sql``` and ```migrations/...down.sql```

### Setp 13.3 Apply migrations
```
migrate -path migrations -database "postgres://localhost:5432/restapi?sslmode=disable&user=postgres&password=postgres" up
```

