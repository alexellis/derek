
# Hacking on Derek

Bear in mind that we offer a managed service, which pays towards the maintenance and upkeep of Derek. This guide should only be used for development of Derek.

Deployment options for development:

* Deploy to faasd (beginner skill level)
* Deploy to k3s (medium skill level)

> You can also `faas-cli build` followed by `docker run` to run Derek's container locally, but this is an advanced option and outweighs the ease of use of deploying faasd or k3s with openfaas.

Summary of what you need to do:

* Set up faasd or openfaas on k3s
* Create a GitHub App
* Install the GitHub App on your repos or orgs
* Create your secrets for derek with `faas-cli secret create`
* Configure `.DEREK.yml`
* Test out Derek

Estimated setup time 30-60mins depending on your experience-level with GitHub's Apps.

## Ready.. Set.. Derek!

### Deploy OpenFaaS

* [Deploy OpenFaaS with faasd](https://github.com/openfaas/faasd)
* [Deploy OpenFaaS with Kubernetes/K3s](https://docs.openfaas.com/deployment/kubernetes/)

### Get your public URL

Now get your publicly-available URL for the OpenFaaS gateway.

If you're behind a firewall, use an [inlets](https://inlets.dev) tunnel - feel free to get a free trial.

```bash
export OPENFAAS_URL=https://gw.example.com
```

### Create a GitHub App in your GitHub account named "Derek dev"

For the homepage URL enter: https://github.com/alexellis/derek

Enter these OAuth Permissions/Scopes:

- Issues - read/write
- Pull requests - read/write
- Repository metadata - read only

If you are setting this up on a private repository you need to grant Derek permissions to download content, this is so that he can download the config file. The write permissions are so that Derek can update your release notes if you are using the `release_notes` feature.

- Repository contents - read/write

Subscribe to these events:

- Issue comment
- Pull request

Set "Where can this GitHub App be installed?" to Any account

* Set your webhook URL to the URL of your OpenFaaS cluster and the `derek` function - i.e. `http://<ip>:<port>/function/derek`

* Note down your Application ID you will need it later where we configure derek.yml with the `github_application_id` variable.

* Click "Generate a private key". This is the private key for your Derek bot, save it as `derek-private-key` and keep it private.

* Set a webhook secret in the UI and save the value as `derek-secret-key`

* Finally install the Derek app onto one of your own GitHub repos. From now on you will get events sent to Derek.

* Create .DEREK.yml in your GitHub repo referring to the docs for what to put in the file.

Next we will build and deploy Derek to your OpenFaaS Cluster and trigger a test event such as `Derek add label: testing`

### Define secrets (Swarm)

Before deploying Derek create two secrets:

* derek-private-key - used for GitHub authz
* derek-secret-key - used for verifying webhooks are from GitHub using HMAC

```bash
$ faas-cli secret create derek-private-key \
   --from-file=./derek-private-key
$ faas-cli secret create derek-secret-key \
  --from-file=./derek-secret-key
```

### Configure Docker image

Most of the configuration is outside of the image, the exception being the version of the OpenFaaS watchdog which is pulled from [GitHub](https://github.com/openfaas/faas/releases); check the desired version is being pulled into the image and go ahead and build.

Update derek.yml and replace `alexellis/` with your own username on the Docker Hub or your registry.

```bash
$ faas-cli build -f derek.yml
```

This will build a Docker image in your local library, now push it to the Docker Hub.

```bash
$ docker login
$ faas-cli push -f derek.yml
```

### Configure stack.yml:

This is where Derek finds the details he needs to do the work he does.  The main areas that will need to be updated are the `application_id` and `customer_url` variables.  The gateway value may also need amending if the gateway is remote.

```yaml
provider:
  name: openfaas
  
functions:
  derek:
    handler: ./derek
    image: derek
    lang: Dockerfile
    environment:
      application_id: <github_application_id>
      customer_url: <your_customers_file_url>
      validate_customers: true
      validate_hmac: true
      debug: true
      write_debug: false
      secret_path: /var/openfaas/secrets/
    secrets:
      - derek-secret-key
      - derek-private-key
```

Fill out the `application_id` variable with the ID of the registered Derek GitHub App.

Provide an `https` location of your `.CUSTOMERS` file.  If hosted on GitHub then this should be location obtained from raw URL by clicking on the file in the UI then clicking "Raw". The default value is: `https://raw.githubusercontent.com/openfaas/faas/master/.DEREK.yml`

Validating via a symmetric key is also known as HMAC. If the webhook secret wasn't set earlier and you want to turn this off (to edit and debug) then set `validate_hmac="false"`

Now deploy Derek:

```bash
$ faas-cli deploy -f derek.yml
```

* Config environmental options

* `application_id` - ID provided in UI for GitHub App
* `customer_url` - A text file of valid repos which can use your Derek installation, separated by new-lines normally `CUSTOMERS`
* `validate_customers` - If set to false then the `customer_url` is ignored
* `validate_hmac` - Validate all incoming webhooks are signed with the secret `derek-secret-key` that you enter in the GitHub UI
* `write_debug` - Dump the incoming request to the function logs. This is not needed since the request can be viewed in the advanced tab of the GitHub App UI

### Configure your first GitHub Repo for Derek

Finally configure the features that you want to enable within your GitHub repo by creating a `.DEREK.yml` file.

The file should detail which features you wish to enable and the maintainer names; for example this repo would look as follows:

```yml
maintainers:
 - alexellis
 - rgee0
 - johnmccabe

features:
 - dco_check
 - comments
```

### Testing and troubleshooting

To test: 

Create a label of "no-dco" within every project you want Derek to help you with.
 
Head over to your GitHub repository and raise a Pull Request from the web-UI for your README file. This will not sign-off the commit, so you'll have Derek on your case.

If you're not sure if things are working right then Click on the GitHub App via your account settings and then click "Advanced" and "Recent Deliveries". This will show you all the incoming messages and their responses.

You can tweak your environment and then hit "Redeliver" to send the message again.

### Appendix

#### Personal Access Tokens

> Note: you may be able to use a Personal Access Token generated through your *Developer Settings* page instead of a GitHub App. If you wish to do this then set the `personal_access_token` environmental variable. It is not recommended since the granularity is much more coarse meaning you must grant more privileges to the bot. The bot will also report as your own GitHub user account.

