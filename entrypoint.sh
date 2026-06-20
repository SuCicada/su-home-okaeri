#!/bin/sh
set -e

RUN_USER="${RUN_USER:-sucicada}"
PUID="${PUID:-1000}"
PGID="${PGID:-1000}"

CURRENT_UID=$(id -u "$RUN_USER" 2>/dev/null || echo -1)
CURRENT_GID=$(getent group "$RUN_USER" | cut -d: -f3 2>/dev/null || echo -1)

if [ "$CURRENT_UID" != "$PUID" ] || [ "$CURRENT_GID" != "$PGID" ]; then
  groupmod -o -g "$PGID" "$RUN_USER"
  usermod -o -u "$PUID" -g "$PGID" "$RUN_USER"
fi

exec runuser -u "$RUN_USER" -- "$@"
