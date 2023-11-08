# Auth0 Actions-As-Code

## Pre-requisites

### Generate the Auth0 management API client credentials

1. Go to the [Auth0 > Applications > Applications](https://manage.auth0.com/dashboard/eu/_/applications) and click *
   *Create Application**.
2. Choose **Machine to Machine Application**.
3. Give it a name like **Actions Updater** and click **Create**.
4. Select the **Auth0 Management API** from the **APIs** list.
5. Check the **create:actions**, **read:actions**, **update:actions** permissions and click **Authorize**.
6. Go to the **Settings** tab and copy the **Domain**, **Client ID** and **Client Secret**.
7. Add the **Domain**, **Client ID** and **Client Secret**
   as [secrets](https://docs.github.com/en/actions/reference/encrypted-secrets) to your repository. Recommended names
   are **AUTH0_TENANT_DOMAIN**, **AUTH0_CLIENT_ID** and **AUTH0_CLIENT_SECRET**.

### Get the Auth0 Action ID

1. Go to the [Auth0 > Actions > Library](https://manage.auth0.com/dashboard/eu/_/actions/library?tab=1) and select the **Custom** tab.
2. Click on the action you want to update.
3. Copy the **Action ID** from the URL. For example in the following action https://manage.auth0.com/dashboard/eu/_/actions/library/details/4cf1a082-ef6f-460c-9ce2-ae6f3b027a68 the id is **4cf1a082-ef6f-460c-9ce2-ae6f3b027a68**.

## Create the config.yml file

1. Basic use case

```yaml
actions:
  post_login:
      - id: '4cf1a082-ef6f-460c-9ce2-ae6f3b027a68'
        name: 'My Very First Post-Login Action'
        code_file_path: './post-login.js'
```

2. Advanced use case

```yaml
actions:
  post_login:
    - id: '4cf1a082-ef6f-460c-9ce2-ae6f3b027a68'
      name: 'My Very First Post-Login Action'
      code_file_path: './post-login.js'
      dependencies:
        - name: 'axios'
        - name: 'lodash'
          version: '1.0.0'
      secrets:
        - key: 'API_BASE'
          value: 'https://api.example.com'
        - key: 'API_TOKEN'
          env_key: 'SOME_TOKEN_KEY_IN_GITHUB_SECRETS'
```

## Inputs

| Name                  | Description                                               | Required | Default    |
|-----------------------|-----------------------------------------------------------|----------|------------|
| `auth0_client_id`     | The Auth0 Client ID.                                      | **✔️**   |            |
| `auth0_client_secret` | The Auth0 Client Secret.                                  | **✔️**   |            |
| `auth0_tenant_domain` | The Auth0 Tenant Domain.                                  | **✔️**   |            |
| `config_path`         | The path to the Auth0 Actions-As-Code configuration file. |          | config.yml |

## Example usage

```yaml
uses: stefanoschrs/auth0-actions-as-code
with:
   auth0_client_id: ${{ secrets.AUTH0_CLIENT_ID }}
   auth0_client_secret: ${{ secrets.AUTH0_CLIENT_SECRET }}
   auth0_tenant_domain: ${{ secrets.AUTH0_TENANT_DOMAIN }}
```

```yaml
uses: stefanoschrs/auth0-actions-as-code
with:
  auth0_client_id: ${{ secrets.AUTH0_CLIENT_ID }}
  auth0_client_secret: ${{ secrets.AUTH0_CLIENT_SECRET }}
  auth0_tenant_domain: ${{ secrets.AUTH0_TENANT_DOMAIN }}
  config_path: ./path-to-config.yml
```
