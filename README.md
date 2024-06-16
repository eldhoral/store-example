# Store Service

Store API Service

## Prerequisites

**Install Go v 1.17**

Please check the [Official Golang Documentation](https://golang.org/doc/install) for installation.

**Install SQL-Migrate**

```bash
go get -v github.com/rubenv/sql-migrate/...
```

**Install Mockery**

```bash
go get github.com/vektra/mockery/v2/.../
```

**Upgrade package**

```
- Upgrade single package
go get -u github.com/gorilla/mux
go mod tidy

- Upgrade all
go get -u
go mod tidy

```


## Installation


**Download dependencies (optional)**

If you want to download all dependencies into the vendor folder, please run the following command:

```bash
go mod vendor
```

**Clone this repository**

```bash
git clone store-api.git
# Switch to the repository folder
cd paylater-customer-api
```

**Copy the `.env.example` to `.env`**

```bash
cp .env.example .env
```

Make the required configuration changes in the `.env` file.

**Copy the `dbconfig.yml.example` to `dbconfig.yml`**

```bash
cp dbconfig.yml.example dbconfig.yml
```

Make the required configuration changes in the `dbconfig.yml` file.

**Run DB Migration**

```bash
make migrate-sql
make migrate-data
```

**Run Application**

```bash
make run
```

## Install on Docker

**Create docker image**

```bash
make image
```

## Install on Docker Compose

**Build**

```bash
make docker-staging-build
```

**Run**

```bash
make docker-staging-run
```

## Unit Testing

**Mocking The Interface**
```bash
cd internal/{function folder}
# Mock Repository interface
mockery --name=Repository --output=../mocks
# Mock Service interface
mockery --name=Service --output=../mocks
```

**Run Unit Test**

To run unit testing, just run the command below:
```bash
make test
```

**Code Coverage**

If you want to see code coverage in an HTML presentation (after the test) just run:

```bash
make coverage
```

## Folders

* `cmd` - Contains command files.
* `app/api` - Contains http server.
* `app/docker` - Contains Dockerfile.
* `app/migrations` - Contains DB migrator.
* `internal` - Contains packages which are specific to your project.
* `pkg` - Contains extra packages.

## Reference

* [Folder Explanation](https://github.com/golang-standards/project-layout)
* [Go Modules](https://blog.golang.org/using-go-modules)
* [Google JSON Style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)
* [Gorilla Mux](https://www.gorillatoolkit.org/pkg/mux)
* [Logrus](https://github.com/sirupsen/logrus)
* [Mockery](https://github.com/vektra/mockery)
* [SQL-Migrate](https://github.com/rubenv/sql-migrate)
* [SQLMock](https://github.com/DATA-DOG/go-sqlmock)
* [Testify](https://github.com/stretchr/testify)

## Contributing

When contributing to this repository, please note we have a code standards, please follow it in all your interactions with the project.

#### Steps to contribute

1. Clone this repository.
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Submit pull request.

**Note** :

* Please make sure to update tests as appropriate.

* It's recommended to run `make test` command before submit a pull request.

* Please update the postman collection if you modify or create new endpoint.
* CMS
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/16666203-9e865a79-c825-49c6-ba81-db61d023b1bc?action=collection%2Ffork&collection-url=entityId%3D16666203-9e865a79-c825-49c6-ba81-db61d023b1bc%26entityType%3Dcollection%26workspaceId%3D94b64e98-ca7a-420a-824f-45deac156947)
* Mobile
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/16666203-94ba24a7-a997-4374-bf78-67a3ae7f129f?action=collection%2Ffork&collection-url=entityId%3D16666203-94ba24a7-a997-4374-bf78-67a3ae7f129f%26entityType%3Dcollection%26workspaceId%3D94b64e98-ca7a-420a-824f-45deac156947)

