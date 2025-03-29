## App migrations

- Ensure your control plane localdev is up and running
- Ensure you have [go migrate](https://formulae.brew.sh/formula/golang-migrate) installed on your system.
- Create a migration in this directory by running `migrate create -ext sql <migration_name>` in this directory. Choose a meaningful migration name in snake case -- refer to existing migration names in this folder.
- This would create two files like `20240823093944_migration_name.up.sql` and `20240823093944_migration_name.down.sql`.
- Add an up and down migration SQL in these files.
- Rebuild with `docker compose up -d --build app-migrate`. Your migrations should get applied automatically. They will also get applied on stage and production as part of deployment pipeline.