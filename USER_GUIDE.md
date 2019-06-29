## Derek User Guide

This is a user-guide for all the features available in Derek.

### Get Derek

You can self-host Derek or use the managed service, see [GET.md](./GET.md).

### Enable Derek for a repo

Add a .DEREK.yml file to the repository.

#### .DEREK.yml schema

Example from this project:

```yaml
maintainers:
 - rgee0
 - alexellis
 
features:
 - dco_check
 - comments
 - pr_description_required
```

This file enables Derek usage for `rgee0` and `alexellis`, it also turns on all features available. If you specifically do not want the commenting or `dco_check` feature then comment out the line or remove it from your file. At least one feature is required for Derek to be of use.

### Feature: `dco_check`

If `dco_check` is specified in the feature list then Derek will inform you when a PR is submitted with commits which have no sign-off message. This is required if your project requires the [Developer Certificate of Origin](https://developercertificate.org).

Derek will add a label of `no-dco` and a comment to help the PR submitter fix the commits.

### Feature: `pr_description_required`

If `pr_description_required` is specified in the feature list then Derek will inform you that a PR needs a description. He also adds the `invalid` label.

### Feature: `redirect` config

The .DEREK.yml file can be redirected to another repository or site. This is used in the OpenFaaS project where around 12 repos are present with the same permissions, features and users.

Example:

```yaml
redirect: https://raw.githubusercontent.com/openfaas/faas/master/.DEREK.yml
```

### Feature: `comments`

If `comments` is given in the `features` list then this enables all commenting features as below.

> Note: All commands can be given with a prefix of either `Derek <command>` or `/<command>`.

#### Edit title

Let's say a user raised an issue with the title `I can't get it to work on my computer`

```
Derek set title: Question - does this work on Windows 10?
```
or
```
Derek edit title: Question - does this work on Windows 10?
```

#### Manage labels

Labels can be used to triage work or help sort it.

```
Derek add label: proposal
```
```
Derek add label: help wanted
```
```
Derek remove label: bug
```

To address multiple labels through a single action use a comma separated list.  The maximum number of labels that can be managed in one comment defaults to 5; this can be set to preference through `multilabel_limit` in your `stack.yml`

To add multiple labels:
```
Derek add label: proposal, help wanted, skill/intermediate
```

To remove multiple labels:
```
Derek remove label: proposal, help wanted
```

#### Set a milestone for an issue or PR

You can organize your issues in groups through existing milestones

```
Derek set milestone: example
```
```
Derek remove milestone: example
```

#### Assign work

You can assign work to people too

```
Derek assign: alexellis
```

Use the `me` moniker to refer to yourself

```
Derek unassign: me
```

#### Add a reviewer to a PR

You can assign people for a PR review as well

```
Derek set reviewer: alexellis
```

Use the `me` moniker to refer to yourself

```
Derek clear reviewer: me
```

> Note: Both assigning work and/or PR reviewer rely on the target user being a member of your GitHub organisation or for a personal project, they must be a collaborator with write-access.

#### Open and close issues and PRs

Sometimes you may want to close or re-open issues or Pull Requests:

```
Derek close
```
```
Derek reopen
```

A reason can also be added if further explanation is appropriate:

```
Derek close: not an issue
```
```
Derek reopen: work incomplete
```

#### Lock/un-lock conversation/threads

This is useful for when conversations are going off topic or an old thread receives a lot of comments that are better placed in a new issue. 

```
Derek lock
```
```
Derek unlock
```

> Note: once locked no further comments are allowed apart from users with admin access.

#### Add predefined message

Have Derek add pre-configured comments to a PR or Issue thread, for example when you would like to direct someone towards the contributing guide.

Configure the feature in `.DEREK.yml` file. It should look something like:

```
custom_messages:
  - name: docs
    value: Hello, please check out the docs ...
  - name: slack
    value: |
           -- 
           To join our slack channel ...
```

Above are two examples which shows simple configuration, the first one is the method for single line messages, the second one is more specific multi line literal, which should be exactly below the `|` sign in order to be displayed and not having errors while parsing.

Tell derek to send the message:

```
Derek message: docs
```
```
Derek msg: slack
```

### Notes on usage

#### Editing the .DEREK.yml file

The .DEREK.yml file is served by a GitHub CDN which has a 5 minute cache expiry. That means if you make a change, it will take at least 5 minutes before it kicks in.

#### Multiple-commands in a comment

Multiple commands in a single comment are not yet supported.

#### Additional white-space

Additional white-space/new-lines in comments are not yet supported

### Enroll users to use Derek with your repo

Users can be specified in a list called `curators` or `maintainers` - both are offered for when the term `maintainer` is a loaded term. The Moby project use this variant, for everyone else `maintainers` may make sense.

Usernames are strictly case-sensitive.
