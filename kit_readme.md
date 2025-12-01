#1 Requirements

1. Already installed Go language
2. Existing PostgreSQL database
3. Installed Air (https://github.com/air-verse/air)

#2 Preparation

0. Replace entire repository `github.com/garrychrstn/kit-go` with your own module name.
1. Install sqlc (snap)
`sudo snap install sqlc`

Visit https://docs.sqlc.dev/en/stable/overview/install.html for more information.

2. Install `mgr.go` dependency
a. Cut / Remove `//go:build ignore` from `mgr.go`. 
b. Run `go mod tidy` to install existing dependencies.
c. Put `//go:build ignore` back to `mgr.go`.

3. Create migration using `mgr.go` tools.
Read migration_readme.md for `mgr.go` usage.

4. Generate database.
`sqlc generate`

5. Run using `air`
