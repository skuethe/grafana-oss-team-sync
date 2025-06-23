<!--
SPDX-FileCopyrightText: 2020 contributors to the cncf/project-template project <https://github.com/cncf/project-template>
SPDX-License-Identifier: Apache-2.0
-->

# Contributing Guide

* [New Contributor Guide](#contributing-guide)
  * [Ways to Contribute](#ways-to-contribute)
  * [Find an Issue](#find-an-issue)
  * [Ask for Help](#ask-for-help)
  * [Pull Request Lifecycle](#pull-request-lifecycle)
  * [Development Environment Setup](#development-environment-setup)
  * [Sign Your Commits](#sign-your-commits)
  * [Pull Request Checklist](#pull-request-checklist)

Welcome! We are glad that you want to contribute to our project! 💖

As you get started, you are in the best position to give us feedback on areas of
our project that we need help with including:

* Problems found during setting up a new developer environment
* Gaps in our Quickstart Guide or documentation
* Bugs in our automation scripts

If anything doesn't make sense, or doesn't work when you run it, please open a
bug report and let us know!

## Ways to Contribute

We welcome many different types of contributions including:

* New features
* Builds, CI/CD
* Bug fixes
* Documentation
* Issue Triage

## Find an Issue

We have good first issues for new contributors and help wanted issues suitable
for any contributor. [good first issue][goodfirstissue] has extra information to
help you make your first contribution. [help wanted][helpwanted] are issues
suitable for someone who isn't a core maintainer and is good to move onto after
your first pull request.

Once you see an issue that you'd like to work on, please post a comment saying
that you want to work on it. Something like "I want to work on this" is fine.

## Ask for Help

The best way to reach us with a question when contributing is to ask on:

* The original github issue

## Pull Request Lifecycle

⚠️ **Explain your pull request process**

## Development Environment Setup

⚠️ **Explain how to set up a development environment**

Current tool-set:
- (optional) `podman` - to run lint on `go`, `reuse` as well as run integration and e2e tests locally
- (optional) `podman-compose` - to run integration and e2e tests locally
- (optional) `pre-commit` - get better feedback before you commit your work

1. fork the repo
2. clone it locally
3. initialize pre-commit (`pre-commit install`)
4. (optionally): install `reuse`
- use podman + podman-compose

## Sign Your Commits

### DCO
Licensing is important to open source projects. It provides some assurances that
the software will continue to be available based under the terms that the
author(s) desired. We require that contributors sign off on commits submitted to
our project's repositories. The [Developer Certificate of Origin
(DCO)][dco] is a way to certify that you wrote and
have the right to contribute the code you are submitting to the project.

You sign-off by adding the following to your commit messages. Your sign-off must
match the git user and email associated with the commit.

    This is my commit message

    Signed-off-by: Your Name <your.name@example.com>

Git has a `-s` command line option to do this automatically:

    git commit -s -m 'This is my commit message'

If you forgot to do this and have not yet pushed your changes to the remote
repository, you can amend your commit with the sign-off by running

    git commit --amend -s

## Pull Request Checklist

When you submit your pull request, or you push new commits to it, our automated
systems will run some checks on your new code. We require that your pull request
passes these checks, but we also have more criteria than just that before we can
accept and merge it. We recommend that you check the following things locally
before you submit your code:

  - add the `SPDX` headers to newly created files (with the exception of files in the `test/**` folder)  
    This can be achieved by the [`reuse` binary][reuseinstall]:
    ```
    reuse annotate --copyright="Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>" --license="GPL-3.0-or-later" [PATH TO NEW FILE HERE]
    ```
  - execute `REUSE lint` on the code  
    There is a helper available: you can run `scripts/container.sh licenses` to use podman to run REUSE's `lint` command on the code
  - execute `golangci-lint` on the code  
    There is a helper available: you can run `scripts/container.sh lint` to use podman to run golangci-lint's `run` command on the code


[goodfirstissue]:   <https://github.com/skuethe/grafana-oss-team-sync/labels/good%20first%20issue> "Issues with label 'good first issue'"
[helpwanted]:       <https://github.com/skuethe/grafana-oss-team-sync/labels/help%20wanted> "Issues with label 'help wanted'"
[dco]:              <https://probot.github.io/apps/dco/> "Developer Certificate of Origin"
[reuseinstall]:     <https://github.com/fsfe/reuse-tool?tab=readme-ov-file#install> "How to install REUSE tool"
