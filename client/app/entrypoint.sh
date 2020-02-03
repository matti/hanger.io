#!/usr/bin/env bash
set -euo pipefail

_term() {
  >&2 echo "TERM"
  exit 0
}
trap "_term" TERM

_err() {
  >&2 echo "err: $*"
  exit 1
}

while ! nc -w 1 -z server 8080; do
  echo "waiting for server:8080"
  sleep 1 &
  wait $!
done
echo "server:8080 ready"

while true; do
  case ${1:-} in
    pauser)
      echo "pause"
      curl --silent server:8080/pause/client &
      wait $!
      echo "continue"
    ;;
    continuer)
      echo "continue in 1s"
      sleep 1 &
      wait $!
      curl --silent server:8080/continue/client &
      wait $!
      echo "go"
    ;;
    *)
      _err "unknown arg"
    ;;
  esac
done
