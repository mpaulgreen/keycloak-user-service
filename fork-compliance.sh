#!/bin/bash -ex

# forces a merge request author to add @devtools-bot to their fork.
# this will allow the gitlab-notifier plugin to add pipelines information to the MR,
# which will enable automatic merges by the bot.
# this will also add all members of the group mentioned as the last argument
# to be maintainers of the fork as well.
# current group is appeng-keycloak: https://gitlab.cee.redhat.com/appeng-keycloak

set -exvo pipefail

mkdir -p /tmp/config
echo "$APP_INTERFACE_CONFIG_TOML" | base64 -d > /tmp/config/config.toml

docker run --rm \
    --volume /tmp/config:/config:z \
    --workdir / \
    quay.io/app-sre/qontract-reconcile:latest \
    qontract-reconcile --config /config/config.toml \
    gitlab-fork-compliance \
    $gitlabMergeRequestTargetProjectId $gitlabMergeRequestIid appeng-keycloak
