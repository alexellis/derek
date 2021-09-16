## Onboarding guide for Derek

We host Derek for you, to save you time and effort.

### How do I start using Derek?

If you're considering using Derek, then please use our managed service, available via GitHub Sponsors. Just Sponsor OpenFaaS and you can request access immediately after by sending a pull request to the [CUSTOMERS file](https://github.com/alexellis/derek/blob/master/.CUSTOMERS).

The price covers hosting, management and maintenance of the service.

* [Sponsor OpenFaaS on Github](https://github.com/sponsors/openfaas) on the 25 USD/mo tier

### The workflow

GitHub sends webhooks to Derek for different events and comments that happen across your repositories. Derek then looks for a .DEREK.yml file to see if the repository requires any response.

You can also use a single main repository and then redirect to that from the others, or use a different .DEREK.yml file in each. 

For example:

* [openfaas/faas](https://github.com/openfaas/faas/blob/master/.DEREK.yml) is the main file
* [openfaas/faas-cli](https://github.com/openfaas/faas-cli/blob/master/.DEREK.yml) is a redirect file.

Installation process:

* Install the managed or your self-hosted Derek GitHub App
* Send a PR to the [CUSTOMERS file](https://github.com/alexellis/derek/blob/master/.CUSTOMERS) with your GitHub username or GitHub organization
* Add a .DEREK.yml to any repositories you want to include, turn on or off any features you need as per the [user guide](./USER_GUIDE.md)
* Add in more repositories in the same organisation using the redirect feature

### Installing Derek once you're a customer

To use our managed Derek service follow the instructions below which take around 5-10 minutes.

* Install this GitHub App on the individual repository (not the whole org):

https://github.com/apps/derek

You will be told what permissions are required.

* Create `.DEREK.yml` in your selected repositories. You can use a redirect file if you have several repositories.

You can use this file as a template: https://github.com/openfaas/faas/blob/master/.DEREK.yml

* Now raise a PR to the `.CUSTOMERS` file

Raise a PR to this file, and make sure you use `git commit --signoff` rather than the UI to make the PR

https://github.com/alexellis/derek/blob/master/.CUSTOMERS

* Finally test it works

Raise a new issue and type in `Derek close` or `/close`, then edit your `.DEREK.yml` file to add your team and community maintainers/contributors.

## Contributions

If you're looking to hack on Derek, see [DEV.md](DEV.md) for how to set it up locally for testing and development.

