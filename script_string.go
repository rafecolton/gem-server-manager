package gsm

const SCRIPT_STRING = `#!/usr/bin/env bash

usage() {
  cat <<EOB
  Usage: retrieve-gems

  Downloads the Gemfile and Gemfile.lock for the specified repo and rev,
  retrieves the gems, places them in the specified GEMDIR.

  Require Environmental Variables:
  GEMINABOX_HOST - the host to use for uploading gemiabox
  GEMDIR         - the directory in which to place the downloaded gems
  OWNER          - the GitHub organization for the repo
  REPO           - the repo name
  REV            - the rev at which to download the Gemfile & Gemfile.lock
  AUTH_TOKEN     - the GitHub authorization token
EOB
}

main() {
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
  esac
  validate_env
  download_gemfile
  bundle_exec
  put_things_inabox
}

validate_env() {
  missing=0
  for var in GEMINABOX_HOST GEMDIR OWNER REPO REV AUTH_TOKEN ; do
    if [[ -z "$(eval "echo \$$var")" ]] ; then
      echo "\$$var needs to be defined"
      missing=1
    fi
  done

  if [[ $missing -gt 0 ]] ; then
    usage
    exit 3
  fi
}

download_gemfile() {
  for file in 'Gemfile' 'Gemfile.lock' ; do
    result="$(_get_file "$file")"
    if echo "$results" | json -a message | grep -q -i 'not found' ; then
      echo "Could not locate $file"
      exit 5
    else
      echo "$results" | json -a content | base64 --decode > "$TMP_DIR/$file"
    fi
  done
}

bundle_exec() {
  pushd "$TMP_DIR" >/dev/null
  bundle install --path "$GEMDIR"
  popd
}

put_things_inabox() {
  pushd "$TMP_DIR" >/dev/null
  for filename in $(bundle list | tr "()" " " | \
	awk '{ if (NR > 1) { print "*"$2"-"$3".gem" }}' | \
	xargs -I {} echo "$(find $GEMDIR -type f -name '*.gem' | head -n 1 | xargs dirname)/{}" \
  ); do
	[[ -s "$filename" ]] && gem inabox "$filename"
  done
}

_get_file() {
  file_path="$1"
  curl -s -XGET -H "Authorization: token $AUTH_TOKEN" \
    "https://api.github.com/repos/$OWNER/$REPO/contents/$file_path"
}

export TMP_DIR="$(mktemp -d -t 'XXXXXXXXXX')"
trap "rm -rf $TMP_DIR" EXIT SIGINT SIGTERM
main "$@"`
