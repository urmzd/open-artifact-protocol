<oap:target id="readme-document">
# <oap:target id="project-name">dbmigrate</oap:target>

<oap:target id="overview-section">
![<oap:target id="build-badge-alt">Build</oap:target>](<oap:target id="build-badge-url">https://img.shields.io/github/actions/workflow/status/example/dbmigrate/ci.yml?branch=main</oap:target>)
![<oap:target id="version-badge-alt">Version</oap:target>](<oap:target id="version-badge-url">https://img.shields.io/npm/v/dbmigrate.svg</oap:target>)
![<oap:target id="license-badge-alt">License</oap:target>](<oap:target id="license-badge-url">https://img.shields.io/npm/l/dbmigrate.svg</oap:target>)

<oap:target id="brief-description">dbmigrate is a fast, opinionated database migration CLI for managing schema changes, running migrations safely, and keeping environments in sync.</oap:target>
</oap:target>

<oap:target id="installation-section">
## Installation

<oap:target id="installation-npm-section">
### npm
<oap:target id="installation-npm-command">npm install -g dbmigrate</oap:target>
</oap:target>

<oap:target id="installation-brew-section">
### Homebrew
<oap:target id="installation-brew-command">brew install dbmigrate</oap:target>
</oap:target>

<oap:target id="installation-cargo-section">
### Cargo
<oap:target id="installation-cargo-command">cargo install dbmigrate</oap:target>
</oap:target>
</oap:target>

<oap:target id="usage-section">
## Usage

<oap:target id="usage-intro">Run `dbmigrate --help` to see all available commands and options.</oap:target>

<oap:target id="usage-examples-list">
- <oap:target id="usage-example-init">dbmigrate init</oap:target>
- <oap:target id="usage-example-create">dbmigrate create add_users_table</oap:target>
- <oap:target id="usage-example-create-timestamped">dbmigrate create add_email_index --timestamp</oap:target>
- <oap:target id="usage-example-up">dbmigrate up</oap:target>
- <oap:target id="usage-example-up-specific">dbmigrate up --to 20240501120000</oap:target>
- <oap:target id="usage-example-down">dbmigrate down</oap:target>
- <oap:target id="usage-example-down-count">dbmigrate down --steps 2</oap:target>
- <oap:target id="usage-example-status">dbmigrate status</oap:target>
- <oap:target id="usage-example-rollback">dbmigrate rollback</oap:target>
- <oap:target id="usage-example-rollback-to">dbmigrate rollback --to 20240401093000</oap:target>
- <oap:target id="usage-example-seed">dbmigrate seed</oap:target>
- <oap:target id="usage-example-seed-specific">dbmigrate seed --file seeds/demo.sql</oap:target>
- <oap:target id="usage-example-config">dbmigrate up --config ./dbmigrate.config.json</oap:target>
</oap:target>

<oap:target id="usage-note">Use `dbmigrate status` to verify pending, applied, and failed migrations before deploying changes.</oap:target>
</oap:target>

<oap:target id="configuration-section">
## Configuration

<oap:target id="configuration-file-section">
### Config file format

<oap:target id="configuration-file-description">`dbmigrate` supports JSON, YAML, and TOML configuration files.</oap:target>

<oap:target id="configuration-json-example">
```json
{
  "dialect": "postgres",
  "url": "postgres://user:pass@localhost:5432/app_db",
  "migrationsDir": "./migrations",
  "seedsDir": "./seeds",
  "tableName": "schema_migrations",
  "ssl": false
}
```
</oap:target>

<oap:target id="configuration-yaml-example">
```yaml
dialect: postgres
url: postgres://user:pass@localhost:5432/app_db
migrationsDir: ./migrations
seedsDir: ./seeds
tableName: schema_migrations
ssl: false
```
</oap:target>
</oap:target>

<oap:target id="configuration-env-section">
### Environment variables

<oap:target id="configuration-env-list">
- <oap:target id="env-dbmigrate-url">DBMIGRATE_URL</oap:target>: <oap:target id="env-dbmigrate-url-description">database connection string</oap:target>
- <oap:target id="env-dbmigrate-dialect">DBMIGRATE_DIALECT</oap:target>: <oap:target id="env-dbmigrate-dialect-description">database dialect such as `postgres`, `mysql`, or `sqlite`</oap:target>
- <oap:target id="env-dbmigrate-migrations-dir">DBMIGRATE_MIGRATIONS_DIR</oap:target>: <oap:target id="env-dbmigrate-migrations-dir-description">path to migration files</oap:target>
- <oap:target id="env-dbmigrate-seeds-dir">DBMIGRATE_SEEDS_DIR</oap:target>: <oap:target id="env-dbmigrate-seeds-dir-description">path to seed files</oap:target>
- <oap:target id="env-dbmigrate-table-name">DBMIGRATE_TABLE_NAME</oap:target>: <oap:target id="env-dbmigrate-table-name-description">migration tracking table name</oap:target>
- <oap:target id="env-dbmigrate-log-level">DBMIGRATE_LOG_LEVEL</oap:target>: <oap:target id="env-dbmigrate-log-level-description">log verbosity such as `info`, `debug`, or `error`</oap:target>
</oap:target>
</oap:target>
</oap:target>

<oap:target id="api-section">
## API / Library Usage

<oap:target id="api-intro">You can also use `dbmigrate` as a library from Node.js applications to run migrations programmatically.</oap:target>

<oap:target id="api-example">
```js
import { migrate, rollback, status, createMigration } from 'dbmigrate';

const result = await migrate({
  url: process.env.DBMIGRATE_URL,
  migrationsDir: './migrations'
});

console.log(result.applied);

await status({
  url: process.env.DBMIGRATE_URL
});

await rollback({
  url: process.env.DBMIGRATE_URL,
  steps: 1
});

await createMigration({
  name: 'add_audit_columns',
  migrationsDir: './migrations'
});
```
</oap:target>

<oap:target id="api-methods-list">
- <oap:target id="api-method-migrate">`migrate(options)`</oap:target>: <oap:target id="api-method-migrate-description">applies pending migrations</oap:target>
- <oap:target id="api-method-rollback">`rollback(options)`</oap:target>: <oap:target id="api-method-rollback-description">reverts the last applied migration or a specified number of steps</oap:target>
- <oap:target id="api-method-status">`status(options)`</oap:target>: <oap:target id="api-method-status-description">returns applied, pending, and failed migration information</oap:target>
- <oap:target id="api-method-create-migration">`createMigration(options)`</oap:target>: <oap:target id="api-method-create-migration-description">creates a new migration file</oap:target>
</oap:target>
## Troubleshooting

- **`dbmigrate` cannot connect to the database**: Verify `DBMIGRATE_URL` or your config file, confirm the database is running, and check network/firewall settings.
- **Migrations fail with a syntax error**: Inspect the generated migration file for invalid SQL or schema changes that are not supported by your database dialect.
- **A migration was applied but not recorded correctly**: Check the migrations tracking table name in your config and ensure the application user has permission to read/write it.
- **Seed files are not being found**: Confirm `DBMIGRATE_SEEDS_DIR` or `seedsDir` points to the correct folder and that the file path used in the command is valid.
- **Rollback stops partway through**: Make sure each migration has a valid down/rollback path, then retry with a smaller number of steps or a specific target version.
</oap:target>

<oap:target id="contributing-section">
## Contributing

<oap:target id="contributing-summary">Contributions are welcome. Please open an issue for discussion before making larger changes, follow the existing code style, add tests for new behavior, and ensure the test suite passes before submitting a pull request.</oap:target>

<oap:target id="contributing-checklist">
- <oap:target id="contributing-fork">Fork the repository</oap:target>
- <oap:target id="contributing-branch">Create a feature branch</oap:target>
- <oap:target id="contributing-tests">Add or update tests</oap:target>
- <oap:target id="contributing-lint">Run linting and formatting</oap:target>
- <oap:target id="contributing-pr">Open a pull request</oap:target>
</oap:target>
</oap:target>

<oap:target id="license-section">
## License

<oap:target id="license-name">MIT</oap:target>
</oap:target>
</oap:target>