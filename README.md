# go-dal
#### Utility and implement base `database_helper` interface to support data access layer.
#### Init and maintain by LamTNB (baolam0307@gmail.com)
## Installation
Latest release: 0.1.0

### feature
- Create PostgreSQL connection with ORM mapping
- Support mapping field and column by `tag name`
- Support CRUD and some common function
- Support UUID
- Support series ID
- Support Versioning Object

### Getting Started
#### Run unit test
1. Checkout https://github.com/LTNB/go-dal.git
2. Change database connection in `go-dal/postgres/postgres_test.go`:
```
     conf := goDal.Config{
		DriverName:     "postgres",
		DataSourceName: "postgres://postgres:123456@localhost:5432/template?sslmode=disable&client_encoding=UTF-8",
		MaxOpenConns:   5,
		MaxLifeTime:    1 * time.Minute,
		MaxIdleConns:   5,
	}
```

3. Create the table:
```
create table account
(
	id serial not null
		constraint account_pk
			primary key,
	date timestamp,
	email varchar(255),
	full_name varchar,
	role varchar,
	active boolean,
	version integer default 0 not null
);
```
3. Run ```go test -v ./...```

#### Project import
1. Import ```github.com/LTNB/go-dal```
2. Create the configuration:
    ```
   db := go_dal.Config{
            DriverName: driverName,
       		DataSourceName: dataSourceName,
       		MaxIdleConns:   maxIdleConns,
            MaxOpenConns: maxOpenConns,
            MaxLifeTime: maxLifeTime
        }
   db.Init()
   ```
3. Create entity:
```
type AccountMock struct {
	helper.BaseBo  `promoted:"true" id:"uuid"`
	helper.Version `promoted:"true"`
	helper.Auditor        `promoted:"true"`
	Email          string `json:"email" sql:"email"`
	FullName       string `json:"full_name" sql:"full_name"`
	Role           string `json:"role" sql:"role"`
	Active         bool   `json:"active" sql:"active"`
}
``` 
- At `BaseBo`: declare `id` with type `uuid`, `series` or `timestamp` for primary key type
- Including the `Version` if using versioning object
- Including the `Auditor` if using audit log 
4. Examples use ORM mapping: https://github.com/baolam0307/go-dal/blob/master/postgres/postgres_test.go

5. Examples use query builder: https://github.com/baolam0307/go-dal/blob/master/helper/sql/query_builder_test.go

## Usage
`import database/sql`

`github.com/lib/pq`

`import github.com/google/uuid`
