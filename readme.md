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


## Working with migrations

### Step 14. Revert migration
To execute revert ```migrate -path migrations -database "postgres://localhost:5432/restapi?sslmode=disable&user=postgres&password=postgres" down```


### Шаг 1. Новая миграция
Open file ```migrations/.....up.sql```
```
CREATE TABLE users (
    id bigserial not null primary key,
    login varchar not null unique,
    password varchar not null
);

CREATE TABLE articles (
    id bigserial not null primary key,
    title varchar not null unique,
    author varchar not null,
    content varchar not null
);
```

Execute command ```migrate -path migrations -database "postgres://localhost:5432/restapi?sslmode=disable&user=postgres&password=postgres" down```

### Step 2. Define models
To define models ```internal/app/models/``` 2 models:
* user.go
* article.go

```
//user.go
package models

//User model defeniton
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

```

```
//article.go
package models

//Article model defenition
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

```

### Step 3. Define  "repositories"
Working with models through repositories. To do so initialize 2 files:
* ```storage/userrepository.go```
* ```storage/articlerepository.go```

```
//articlerepository.go
package storage

//Instance of Article repository (model interface)
type ArticleRepository struct {
    storage *Storage
}

```

Alike for users.

### Step 4. Allocating public access to the repository
We want our application to communicate with models through repositories (which will contain the necessary set of methods to interact with the database). We need to define 2 methods at the repository, which will provide public repositories:
```
//storage.go

//Instance of storage
type Storage struct {
	config *Config
	// DataBase FileDescriptor
	db *sql.DB
	//Subfield for repo interfacing (model user)
	userRepository *UserRepository
	//Subfield for repo interfaceing (model article)
	articleRepository *ArticleRepository
}

....

//Public Repo for Article
func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return nil
}

//Public Repo for User
func (s *Storage) Article() *ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}
	s.articleRepository = &ArticleRepository{
		storage: s,
	}
	return nil
}

```

### Step 5. What will UserRepo be able to do?
* Save a new user to the database (INSERT user or Create)
* For authentication, you need a user search function by ```login```.
* Output all users from the database
```
package storage

import (
	"fmt"
	"log"

	"github.com/vlasove/go2/7.ServerAndDB2/internal/app/models"
)

//Instance of User repository (model interface)
type UserRepository struct {
	storage *Storage
}

var (
	tableUser string = "users"
)

//Create User in db
func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING id", tableUser)
	if err := ur.storage.db.QueryRow(query, u.Login, u.Password).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

//Find user by login
func (ur *UserRepository) FindByLogin(login string) (*models.User, bool, error) {
	users, err := ur.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var userFinded *models.User
	for _, u := range users {
		if u.Login == login {
			userFinded = u
			founded = true
			break
		}
	}
	return userFinded, founded, nil
}

//Select all users in db
func (ur *UserRepository) SelectAll() ([]*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableUser)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Prepare, where we going to read
	users := make([]*models.User, 0)
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

```
### Step 6. What is needed from ArticleRepo?
* To be able to add an article to the database
* Be able to delete by id
* Receive all articles
* Retrieve an article by id
* Update (at home)
```
articlerepository.go
package storage

import (
	"fmt"
	"log"

	"github.com/vlasove/go2/7.ServerAndDB2/internal/app/models"
)

//Instance of Article repository (model interface)
type ArticleRepository struct {
	storage *Storage
}

var (
	tableArticle string = "articles"
)

//Add article to DB
func (ar *ArticleRepository) Create(a *models.Article) (*models.Article, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, author, content) VALUES ($1, $2, $3) RETURNING id", tableArticle)
	if err := ar.storage.db.QueryRow(query, a.Title, a.Author, a.Content).Scan(&a.ID); err != nil {
		return nil, err
	}

	return a, nil

}

//Delete article by ID
func (ar *ArticleRepository) DeleteById(id int) (*models.Article, error) {
	article, ok, err := ar.FindArticleById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableArticle)
		_, err := ar.storage.db.Exec(query, id)
		if err != nil {
			return nil, err
		}
	}
	return article, nil
}

//Retrieve article by ID
func (ar *ArticleRepository) FindArticleById(id int) (*models.Article, bool, error) {
	articles, err := ar.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var articleFinded *models.Article
	for _, a := range articles {
		if a.ID == id {
			articleFinded = a
			founded = true
			break
		}
	}
	return articleFinded, founded, nil
}

//Get all articles from DB
func (ar *ArticleRepository) SelectAll() ([]*models.Article, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableArticle)
	rows, err := ar.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Prepare where we are going to read
	articles := make([]*models.Article, 0)
	for rows.Next() {
		a := models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Author, &a.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		articles = append(articles, &a)
	}
	return articles, nil
}

```

### Step 7. Description of the router for this project.
Enter ```api```
```
// Trying to configure the router (specifically the router API field)
func (a *API) configreRouterField() {
	a.router.HandleFunc(prefix+"/articles", a.GetAllArticles).Methods("GET")
	a.router.HandleFunc(prefix+"/articles/{id}", a.GetArticleById).Methods("GET")
	a.router.HandleFunc(prefix+"/articles/{id}", a.DeleteArticleById).Methods("DELETE")
	a.router.HandleFunc(prefix+"/articles", a.PostArticle).Methods("POST")
	a.router.HandleFunc(prefix+"/user/register", a.PostUserRegister).Methods("POST")

}
```

Create file ```internal/app/api/handlers.go```
```
```