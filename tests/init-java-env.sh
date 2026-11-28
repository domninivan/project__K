#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
fixture="$(mktemp -d)"
trap 'rm -rf -- "$fixture"' EXIT
mkdir -p "$fixture/JavaBackend/production-java-backend"
printf '# Example\nJWT_SECRET_BASE64=\n' > "$fixture/JavaBackend/production-java-backend/.env.example"
env_file="$fixture/JavaBackend/production-java-backend/.env"
run() { bash "$root/.devcontainer/init-java-env.sh" "$fixture"; }

printf 'JWT_SECRET_BASE64=existing-legacy-key\nOTHER_SETTING=keep\n' > "$fixture/JavaBackend/.env"
legacy_hash="$(sha256sum "$fixture/JavaBackend/.env" | cut -d ' ' -f 1)"
run
[[ ! -e "$fixture/JavaBackend/.env" ]]
[[ "$legacy_hash" == "$(sha256sum "$env_file" | cut -d ' ' -f 1)" ]]
rm -- "$env_file"

run
[[ "$(sed -n 's/^JWT_SECRET_BASE64=//p' "$env_file" | base64 -d | wc -c)" == 32 ]]
first="$(sha256sum "$env_file")"
run
[[ "$first" == "$(sha256sum "$env_file")" ]] || { echo 'Existing key changed'; exit 1; }
printf 'OTHER_SETTING=keep\r\nJWT_SECRET_BASE64=\r\n' > "$env_file"
run
grep -q '^OTHER_SETTING=keep' "$env_file"
[[ "$(sed -n 's/^JWT_SECRET_BASE64=//p' "$env_file" | base64 -d | wc -c)" == 32 ]]
printf 'OTHER_SETTING=keep' > "$env_file"
run
grep -q '^OTHER_SETTING=keep' "$env_file"
grep -q '^JWT_SECRET_BASE64=' "$env_file"
printf 'JWT_SECRET_BASE64=\047\047\n' > "$env_file"
run
[[ "$(sed -n 's/^JWT_SECRET_BASE64=//p' "$env_file" | base64 -d | wc -c)" == 32 ]]
printf 'JWT_SECRET_BASE64=existing-one\nJWT_SECRET_BASE64=existing-two\n' > "$env_file"
before="$(sha256sum "$env_file")"
if run; then echo 'Duplicate definitions accepted'; exit 1; fi
[[ "$before" == "$(sha256sum "$env_file")" ]]
echo '[OK] Initialization: new file, preserved key, empty value, missing entry, CRLF, duplicate rejection.'
