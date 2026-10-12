#!/bin/bash
set -e

echo "Menjalankan database schema..."

docker compose exec -T postgres \
  psql -v ON_ERROR_STOP=1 \
  -U lebah -d familymart \
  < db/schema.sql

echo "Menjalankan database seed..."

docker compose exec -T postgres \
  psql -v ON_ERROR_STOP=1 \
  -U lebah -d familymart \
  < db/seed.sql

echo "Mendaftarkan database functions..."

docker compose exec -T postgres \
  psql -v ON_ERROR_STOP=1 \
  -U lebah -d familymart \
  < db/functions.sql

echo "Database schema, seed, dan functions berhasil dijalankan."