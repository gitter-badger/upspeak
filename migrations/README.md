# migrate CLI usage

## Installation

```
$ cd $GOPATH/bin
$ curl -L https://github.com/mattes/migrate/releases/download/v3.0.1/migrate.linux-amd64.tar.gz | tar xvz
```


## Usage

- Create file name in format `SEQUENCE_COMMENT.TYPE.EXTENSION` in `migrations` dir

> Ex: 1_initialize_teams.up.sql and 1_initialize_teams.down.sql

Ref: [How to write migrations](https://github.com/mattes/migrate/blob/master/MIGRATIONS.md)

```
$ migrate -help
Usage: migrate OPTIONS COMMAND [arg...]
       migrate [ -version | -help ]

Options:
  -source          Location of the migrations (driver://url)
  -path            Shorthand for -source=file://path
  -database        Run migrations against this database (driver://url)
  -prefetch N      Number of migrations to load in advance before executing (default 10)
  -lock-timeout N  Allow N seconds to acquire database lock (default 15)
  -verbose         Print verbose logging
  -version         Print version
  -help            Print usage

Commands:
  goto V       Migrate to version V
  up [N]       Apply all or N up migrations
  down [N]     Apply all or N down migrations
  drop         Drop everyting inside database
  force V      Set version V but don't run migration (ignores dirty state)
  version      Print current migration version
```


So let's say you want to run the first two migrations

```
$ migrate -database postgres://localhost:5432/database up 2
```

If your migrations are hosted on github

```
$ migrate -source github://mattes:personal-access-token@mattes/migrate_test \
    -database postgres://localhost:5432/database down 2
```

The CLI will gracefully stop at a safe point when SIGINT (ctrl+c) is received.
Send SIGKILL for immediate halt.

## Reading CLI arguments from somewhere else

##### ENV variables

```
$ migrate -database "$MY_MIGRATE_DATABASE"
```
