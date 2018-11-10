## Derek User Guide

This is a user-guide for all the features available in Derek.

### Get Derek

See [Get Derek](./GET.md) instructions

### Enable Derek for a repo

Add a .DEREK.yml file to the repository.

### .DEREK.yml schema

Example from this project:

```yaml
curators:
 - rgee0
 - alexellis
 
features:
 - dco_check
 - comments
```

This file enables Derek usage for `rgee0` and `alexellis`, it also turns on all features available. If you specifically do not want the commenting or `dco_check` feature then comment out the line or remove it from your file. At least one feature is required for Derek to be of use.

### Notes on usage

#### Editing the .DEREK.yml file

The .DEREK.yml file is served by a GitHub CDN which has a 5 minute cache expiry. That means if you make a change, it will take at least 5 minutes before it kicks in.

#### Multiple-commands in a comment

Multiple commands are not supported in a comment

#### Additional white-space

Additional white-space or new-lines is not supported

### Enroll users to use Derek with your repo

Users can be specified in a list called `curators` or `maintainers` - both are offered for when the term `maintainer` is a loaded term. The Moby project use this variant, for everyone else `maintainers` may make sense.

Usernames are strictly case-sensitive.

### Feature: redirect config

The .DEREK.yml file can be redirected to another repository or site. This is used in the OpenFaaS project where around 12 repos are present with the same permissions, features and users.

Example:

```yaml
redirect: https://raw.githubusercontent.com/openfaas/faas/master/.DEREK.yml
```

### Feature: `dco_check`

If `dco_check` is specified in the feature list then Derek will inform you when a PR is submitted with commits which have no sign-off. He also adds a label `no-dco`.

### Feature: `pr_description_required`

If `pr_description_required` is specified in the feature list then Derek will inform you that a PR needs a description. He also adds the `invalid` label.

### Feature: `comments`

If `comments` is given in the `features` list then this enables all commenting features:

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
```
Derek unassign: me
```

> Note: This relies on the target user being a member of your GitHub organisation or for a personal project, they must be a collaborator with write-access.


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

> Note: if you lock an issue as an unprivileged user, then only a repository admin can unlock it.

