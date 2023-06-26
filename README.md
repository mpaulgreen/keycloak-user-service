# UserService

## Run Locally

* Run local instance of Keycloak on port 8080
```shell
    docker run -p 8080:8080 --env KEYCLOAK_ADMIN=admin --env KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak start-dev
```

* Option Run userservice with the make command, which is exposed on port 8000.
```shell
    make run-local
```

* Check if userservice returns a http 200 OK response.
```shell
curl --location 'http://localhost:8000/health/live'
```

## Run and Deploy Service on Openshift
Create the template with below command 

`oc process -f dev-template.yaml | oc create -n fips-test -f -`

## Using FindUser API
### Setup local Keycloak server instance 
 * Use a base realm export [here](https://drive.google.com/file/d/1TC46pKENxoYim-zkdq7UkJREriY5HvOf/view) to initialize the keycloak server

### Use Find Users API in userservice
* Get token
```shell
export TOKEN=$(curl -v -X POST --location 'http://(keycloak host url)/realms/(desired realm)/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'username=admin' \
--data-urlencode 'password=admin' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_id=admin-cli' | jq '.access_token')
```

* Call Find Users API with the following example criteria
```shell
# Search users with usernames
curl --location 'http://localhost:8000/users??offset=0&limit=100&org_id=1010101' \
--header "Authorization: Bearer $TOKEN" \
```

### Use Invite User API in userservice

# Send invite
```shell
curl -v -X PUT --location 'http://localhost:8000/user/invite' \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer $TOKEN" \
--data '{"emails": ["user1@company.com", "user2@company.com"], "isAdmin": true, "orgId": 123}'
```

Prerequisites for successful completion:
* Create a group 'CUSTOMER' with 'organization.id' attribute set to '123'
* Create subgroups "admins" and 'users'

## Docker Tasks using Makefile
* Build userservice Docker image
```shell
make  dev-container-build
```

* Push userservice Docker image
```shell
make dev-container-push
```

## Using environment variables
* Use `local-exec-env-vars.sh` as needed to appropriately configure environment variables for the userservice
