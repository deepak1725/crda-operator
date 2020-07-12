#!/usr/bin/env bash
# This script is used as a liveness probe.

PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -U $POSTGRES_USER -d $POSTGRESQL_INITIAL_DATABASE -tc "SELECT 1;" | grep -q 1
