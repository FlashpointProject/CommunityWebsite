# Flashpoint Community Website

WIP website to cover extra community wants, e.g Playlist sharing, Game of the Day.

# Setup

Requires:
- Node 18+
- Go 1.20+

### Setup Client

Navigate to `./web`

Run `npm install`

### Setup Server

Navigate to `./`

Copy and modify `.env.template` to `.env`

Start database:
- Windows: Run `docker-compose -p fpcomm -f dc-db.yml up -d`
- Linux: Run `make db`

Run migrations:
- Windows: Figure it out yourself from the Makefile
- Linux: Run `make migrate`

Use `make rebuild-postgres` to do a complete database reset

# Building + Running

### Client side

Navigate to `./web`

Run `npm run build` to do a single build or `npm run watch` to continously build changes, files will output to `./web/dist`

### Server side

Navigate to `./`

Then:
- Windows: Run `go run ./main/main.go`
- Linux: Run `make run`

The Go server will automatically serve both sides correctly over the same port