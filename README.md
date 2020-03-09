# derek

[![Build Status](https://travis-ci.com/alexellis/derek.svg?branch=master)](https://travis-ci.com/alexellis/derek)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![](https://godoc.org/github.com/alexellis/derek?status.svg)](https://godoc.org/github.com/alexellis/derek)
[![Derek App](https://alexellis.o6s.io/badge?owner=alexellis&repo=derek)](https://github.com/alexellis/derek/)

![Derek's logo](https://pbs.twimg.com/media/DPo4OyrWsAAOk_i.png)

Derek reduces fatigue for maintainers by automating governance and delegating permissions to your team and community.

> Follow [@derekapp](https://twitter.com/derekapp) on Twitter.

## User guide / documentation

[From the team](https://github.com/alexellis/derek#maintainers-and-contributors) that brought you [OpenFaaS](https://www.openfaas.com) - Serverless Functions Made Simple.

### Core features

All features are modular and can be enabled/disabled as required:

* Generate changelogs for releases with PRs merged and commits added
* Let designated non-admin users manage Issues and PRs by commenting `Derek <command>` or `/command`
* Enforce [Developer Certificate of Origin (DCO) checking (optional)](https://developercertificate.org)
* Automatically label/flag PRs without a *Description*
* Label PRs from first-time contributors
* Detect spam PRs (from hacktoberfest)

Self-host, or use the free managed version.

### User guide

Find out what Derek can do you for your project, community and team including all available commands and configuration options: [Read the user-guide](./USER_GUIDE.md)

### Get Derek

You can use a managed Derek bot or self-host your own on OpenFaaS. Start here: [Get Derek](GET.md).

### How do I work?

I'm designed to be installed as a [GitHub App](https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/), but don't worry - I don't need a lot of permissions. Just access to issues and Pull Requests will do.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so deploy me to a Serverless platform like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

This is me in action with the DCO-checking feature. If I find a commit that isn't signed off then I'll let the contributor know in a comment. Normally they edit their commits and re-push within a few minutes after re-reading the message.

![](https://user-images.githubusercontent.com/6358735/29704343-542a36da-8971-11e7-871e-da30c8e86cae.png)

*Inspiration for Derek*

The idea for a bot that could comment on issues or respond to activity is from the docker/docker or Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login). Derek is much simpler (so hackable) and can be installed with granular permissions.

[Live demo here](https://twitter.com/alexellisuk/status/905694832445804544)

### Our users

Some of our users include:

* [Docker / Moby](https://github.com/moby/moby/issues/35736)

* [OpenFaaS](https://github.com/openfaas/faas-cli/issues/85)

* [inlets](https://github.com/inlets/inlets)

* [arkade](https://github.com/alexellis/arkade)

* [k3sup](https://github.com/alexellis/k3sup)

* [goreleaser](https://github.com/goreleaser/goreleaser/commit/9d418755dd9f37589f5a97f34b2c47e2e2f1325a)

* [evilsocket](https://github.com/evilsocket)

* [Subsurface Diving app](https://github.com/Subsurface-divelog/subsurface/pull/1748)

* [okteto](https://github.com/okteto/okteto)

See other users in our [customers file](./.CUSTOMERS)

## Backlog:

* [x] Derek as a managed GitHub App (SaaS)
* [x] Lock thread
* [x] Edit title of Issue/PR
* [x] Toggle the DCO-feature for those who don't need compliance
* [x] Show help messages such as how to join Slack or asking for a rebase via `Derek message`
* [x] Automate release notes based upon PRs closed
* [x] Add each commit to release notes

Future work:

* [ ] Add merge feature based upon policy (in-progress)
* [ ] Caching of customers / .DEREK.yml file
* [ ] Observability of GitHub API Token rate limit
* [ ] Add roles & actions
* [ ] PR target branch checking (i.e. no PRs from master to master)

## Maintainers and contributors

There are [over a dozen contributors](https://github.com/alexellis/derek/graphs/contributors) to Derek, but the following people have made the most impact.

Maintainers:

* Alex Ellis - creator & lead maintainer
* Martin Dekov ([@martindekov](https://github.com/martindekov)) - co-maintainer

Hall of fame:

* Burton Rheutan ([@burtonr](https://github.com/burtonr))
* John McCabe ([@johnmccabe](https://github.com/johnmccabe))
* Lubomir I. Ivanov ([@neolit123](https://github.com/neolit123))
* Ivana Yovcheva ([@ivanayov](https://github.com/ivanayov))
* Richard Gee ([@rgee0](https://github.com/rgee0))

Alex Ellis created Derek to automate project maintainer duties around licensing and to help bring granular permissions back to GitHub. Derek has empowered contributors in the OpenFaaS community to run and maintain the project without needing full write access. OpenFaaS contributors continue to improve Derek so they can get the job done without fuss.

### Contributions

Please follow the [OpenFaaS contribution guide](https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md).
