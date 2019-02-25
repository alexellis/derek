# derek

[![Build Status](https://travis-ci.org/alexellis/derek.svg?branch=master)](https://travis-ci.org/alexellis/derek)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![](https://godoc.org/github.com/alexellis/derek?status.svg)](https://godoc.org/github.com/alexellis/derek)


It's derek ![](https://pbs.twimg.com/media/DPo4OyrWsAAOk_i.png). Nice to meet you. I'd like to help you with Pull Requests and Issues on your GitHub project.

> Please show support for the project and **Star** the repo.

From the team that brought you [OpenFaaS](https://www.openfaas.com) - Serverless Functions Made Simple.

## Core features

* Issue and PR administration through comments for non-admin users
* [Developer Certificate of Origin (DCO) checking (optional)](https://developercertificate.org)
* Reject PRs without descriptions
* Self-host or use the free, managed service 

## User guide / documentation

Find out what Derek can do you for your project, community and team including all available commands and configuration options.

* Read the [user-guide](./USER_GUIDE.md)

## Our users

Some of our users include:

* Docker / Moby:

https://github.com/moby/moby/issues/35736

* OpenFaaS

https://github.com/openfaas/faas-cli/issues/85

* [Subsurface Diving app](https://subsurface-divelog.org)

Example: https://github.com/Subsurface-divelog/subsurface/pull/1748

## Get your own Derek bot for free

You can either provision your own OpenFaaS cluster and install your own private Derek, or use the shared, managed Derek bot for free.

Setup:

* Install the managed or your self-hosted Derek GitHub App
* Send a PR to the [customers file](https://github.com/alexellis/derek/blob/master/.CUSTOMERS) with your GitHub username or GitHub organization
* Finally add your .DEREK.yml - you can use the file from this repository as an example
* Add any other repos optionally using the redirect feature

Start here: [Get Derek](GET.md).

### Backlog:

* [x] Derek as a managed GitHub App
* [x] Lock thread
* [x] Edit title
* [x] Toggle the DCO-feature

Future work:

* [ ] Caching of customers / .DEREK.yml file
* [ ] Observability of GitHub API Token rate limit
* [ ] Add roles & actions
* [ ] Branch Checking

[Live demo here](https://twitter.com/alexellisuk/status/905694832445804544)

## How do I work?

I'm designed to be installed as a [GitHub App](https://developer.github.com/apps/building-integrations/setting-up-and-registering-github-apps/), but don't worry - I don't need a lot of permissions. Just access to issues and Pull Requests will do.

I'm not a long-running daemon.. I'd get bored that way. I work with webhooks - so stick me in a serverless framework like [OpenFaaS](https://github.com/alexellis/faas) and forget about me. Just apply oil from time to time.

This is me in action! Normally contributors edit and re-push within a few minutes after re-reading the contribution guide.

![](https://user-images.githubusercontent.com/6358735/29704343-542a36da-8971-11e7-871e-da30c8e86cae.png)

*Inspiration for Derek*

The idea for a bot that could comment on issues or respond to activity is from the docker/docker or Moby project's bot called [Poule](https://github.com/icecrime/poule). It's a much more complex long-running daemon which uses Personal Access Tokens (so needs to run as a full GitHub login). Derek is much simpler (so hackable) and can be installed with granular permissions.

### Where is Derek working now?

Derek is active and operating 24/7 helping the award-winning OpenFaaS project!

* https://github.com/moby/moby
* http://github.com/openfaas/faas
* http://github.com/openfaas/faas-netes
* http://github.com/openfaas/faas-cli

### Maintainers / contributors

* Alex Ellis - author
* Richard Gee (@rgee0) - co-maintainer
* John Mccabe (@johnmccabe) - contributor

Alex Ellis created Derek to automate project maintainer duties around licensing and to help bring granular permissions back to GitHub. Derek has empowered contributors in the OpenFaaS community to run and maintain the project without needing full write access. OpenFaaS contributors continue to improve Derek so they can get the job done without fuss.

### Contributions

Please follow the [OpenFaaS contribution guide](https://github.com/openfaas/faas/blob/master/CONTRIBUTING.md).

