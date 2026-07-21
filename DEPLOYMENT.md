**Deployment and rollback guide**

- **Purpose**: Recreate DB schema, build images, and run app with Docker Compose.
- **Warning**: The DB recreate step will drop existing data.

Prerequisites
- Docker and Docker Compose installed on the server.
- MySQL port 3306 available (or change compose file).

Quick start (destructive, local)
1. From the repo root run:
```bash
docker compose up -d --build
```
2. The backend will be available on `http://localhost:8080/` and frontend on `http://localhost:3000/`.

Recreate DB schema (destructive)
1. If you need to force schema recreate (drops data), exec into DB container and run script:
```bash
docker compose up -d db
docker cp backend/db/schema.sql $(docker-compose ps -q db):/schema.sql
docker exec -it $(docker-compose ps -q db) sh -c "mysql -u root -p\"$MYSQL_ROOT_PASSWORD\" -e 'DROP DATABASE IF EXISTS db_sistemasetku; CREATE DATABASE db_sistemasetku CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'"
docker exec -i $(docker-compose ps -q db) sh -c "mysql -u root -p\"$MYSQL_ROOT_PASSWORD\" db_sistemasetku" < backend/db/schema.sql
```

Non-destructive migration (recommended in production)
- Use a proper migration tool (golang-migrate) and write incremental SQL migrations.

Rollback steps
- If you dropped data accidentally, restore from backup. Always take a DB dump before destructive operations:
```bash
docker exec -i $(docker-compose ps -q db) sh -c "mysqldump -u root -p\"$MYSQL_ROOT_PASSWORD\" db_sistemasetku > /tmp/backup.sql"
docker cp $(docker-compose ps -q db):/tmp/backup.sql ./backup.sql
```

Next steps I will do if you confirm:
- Run `docker compose up -d --build` in this repo and verify endpoints.
- Add minimal healthcheck and monitoring hints.
