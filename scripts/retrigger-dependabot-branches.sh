#!/usr/bin/env bash
#
# So: for some reason (probably down to how we run tests and push and not pull_request open/ sync
# in order to avoid double running tests- and probably down to how dependabot runs: I suspect it
# uses the API which doesn't trigger a push ¯\_(ツ)_/¯) dependabot PRs don't run tests.
#
# This little hacketty hack should allow us to manually trigger those jobs by committing/ force-pushing
# under a real user, and via git

set -e

PREFIX="origin/"

for B in $(git branch -r); do
    if [[ "${B}" == origin/dependabot/* ]]; then
        git checkout ${B#"$PREFIX"}
        git commit --amend --no-edit
        git push -f
    fi
done

git checkout main
