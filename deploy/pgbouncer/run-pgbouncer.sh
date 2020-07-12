#!/bin/bash -e


PGBOUNCER_DIR=${PGDATA}
USERLIST_TXT=${PGBOUNCER_DIR}/userlist.txt
PGBOUNCER_INI=${PGBOUNCER_DIR}/pgbouncer.ini

mkdir -p ${PGBOUNCER_DIR}

cat > ${USERLIST_TXT} << EOF
"${POSTGRESQL_USER}" "${POSTGRESQL_PASSWORD}"
EOF

cat > ${PGBOUNCER_INI} << EOF

[pgbouncer]
# https://pgbouncer.github.io/config.html
pool_mode = transaction
listen_addr = *
listen_port = 5432
auth_file = ${USERLIST_TXT}
# TODO: we'll want to change this to "cert" or "hba"
auth_type = md5
admin_users = ${POSTGRESQL_USER}
max_client_conn = 10000
default_pool_size = 100
# pgweb refuses to work with extra_float_digits
ignore_startup_parameters = extra_float_digits
log_connections = 0
log_disconnections = 0
log_pooler_errors = 1
stats_period = 60
EOF

pg_user=$(id -u)
PG_USER_OPT=""

if [ $pg_user == "0" ]; then
    chown pgbouncer:pgbouncer ${USERLIST_TXT}
    chown pgbouncer:pgbouncer ${PGBOUNCER_INI}
    PG_USER_OPT="-u pgbouncer"
fi

#pgbouncer ${PG_USER_OPT} ${PGBOUNCER_INI}

