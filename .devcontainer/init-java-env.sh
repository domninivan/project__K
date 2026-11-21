#!/usr/bin/env bash
set -euo pipefail

# Runs in a short-lived Docker container before Compose reads app.env_file.
# Optional workspace argument is used by isolated tests.
workspace="${1:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
java_dir="$workspace/JavaBackend/production-java-backend"
env_file="$java_dir/.env"
[[ -d "$java_dir" ]] || { echo '[ERROR] JavaBackend directory not found.' >&2; exit 1; }
[[ ! -L "$env_file" ]] || { echo '[ERROR] JavaBackend/production-java-backend/.env must not be a symbolic link.' >&2; exit 1; }

lock="$java_dir/.env.init-lock"
mkdir "$lock" 2>/dev/null || { echo '[ERROR] Another Java .env initialization is running; check .env.init-lock.' >&2; exit 1; }
temp_file=''
cleanup() {
    [[ -z "$temp_file" ]] || rm -f -- "$temp_file"
    rmdir -- "$lock"
}
trap cleanup EXIT

# Preserve keys created before .env was moved into the Maven project.
legacy_env="$workspace/JavaBackend/.env"
if [[ ! -e "$env_file" && -e "$legacy_env" ]]; then
    [[ ! -L "$legacy_env" ]] || { echo '[ERROR] Legacy .env must not be a symbolic link.' >&2; exit 1; }
    mv -- "$legacy_env" "$env_file"
    echo '[OK] Existing Java .env moved into production-java-backend.'
fi

# Never print the key, including when preserving an existing value.
source_file="$env_file"
if [[ ! -e "$source_file" ]]; then
    source_file="$java_dir/.env.example"
fi
[[ -f "$source_file" ]] || { echo '[ERROR] JavaBackend/production-java-backend/.env.example not found.' >&2; exit 1; }
state="$(awk '
    /^[[:space:]]*(export[[:space:]]+)?JWT_SECRET_BASE64[[:space:]]*=/ {
        count++; value=$0; sub(/^[^=]*=/, "", value);
        sub(/^[[:space:]]*/, "", value); sub(/[[:space:]]*$/, "", value);
        empty=(value == "" || value == "\"\"" || value == "\047\047" || value ~ /^#/);
    }
    END { if (count > 1) print "DUPLICATE"; else if (count == 0) print "MISSING";
          else if (empty) print "EMPTY"; else print "EXISTS"; }
' "$source_file")"
[[ "$state" != DUPLICATE ]] || { echo '[ERROR] Duplicate JWT_SECRET_BASE64 entries; file unchanged.' >&2; exit 1; }
if [[ "$state" == EXISTS && "$source_file" == "$env_file" ]]; then
    echo '[OK] Existing Java JWT key preserved.'
    exit 0
fi
if [[ "$state" == EXISTS ]]; then
    echo '[ERROR] The tracked .env.example must not contain a real JWT key.' >&2
    exit 1
fi

umask 077
temp_file="$(mktemp "$java_dir/.env.init.XXXXXX")"
GENERATED_JWT_KEY="$(head -c 32 /dev/urandom | base64 | tr -d '\r\n')"
export GENERATED_JWT_KEY
awk '
    /^[[:space:]]*(export[[:space:]]+)?JWT_SECRET_BASE64[[:space:]]*=/ {
        print "JWT_SECRET_BASE64=" ENVIRON["GENERATED_JWT_KEY"]; found=1; next;
    }
    { print }
    END { if (!found) print "JWT_SECRET_BASE64=" ENVIRON["GENERATED_JWT_KEY"] }
' "$source_file" > "$temp_file"
unset GENERATED_JWT_KEY

# Files created as root in the helper must belong to the workspace owner on Linux.
if [[ "$(id -u)" == 0 ]]; then
    chown --reference="$java_dir" "$temp_file"
fi
chmod 600 "$temp_file"
mv -- "$temp_file" "$env_file"
temp_file=''
echo '[OK] Persistent Java JWT key initialized in ignored JavaBackend/production-java-backend/.env.'
