# Knight-Link

### Final Project for COP4710

This is the current repository for the backend (and possibly frontend for the event webapp that we have to create)

Most of the endpoints created are not working now.

Currently working on making the DB and backend as per the ERD provided.

## Goland + PostgresSQL (docker)

## Dev Setup Instructions

### 1. Env Setup:

- Clone the repository: `git clone https://github.com/bingKegeta/Knight-Link.git`

- `cd` into the repo.

- Create a file `.env` and format it like this:

        PG_USER=
        PG_DB=
        PG_PW=

### 2. Database Setup (Docker):

- Make sure to have `docker` and `docker-compose` installed and set up for use.

- `cd` into the cloned repo.
- Run the command `docker compose up`
- Access the database in another terminal using the command `psql -h localhost -U **<PG_USER>** -d **<PG_DB>**`

Rest in progress...
