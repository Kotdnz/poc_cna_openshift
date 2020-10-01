# This repository contain the source files for POC - new applocation in OpenShift cluster.
## Components:
1. BackEnd golang application. We will use it on a two roles - as a client only and as a Server only. Purpose - emulate production and admin wing.
2. Database for storing our test data MongoDB - OpenShift application launched from the OS Marketplace. Nothing special - create new and delate existed dicuments.
3. FrontEnd - versy simple html pages. One as for admin page - able to create new or delete cna user. Second one - production page to viewonly model.
4. Last part - the folder with command for OpenShift to maniulate with the cluster (oc commands and yaml files)

## Architecture:
Inside one project/namespace we have 3 internal services/applications and the route, who exosing our solution and managing level 7 trafic.

|  |  |  |
| -- | -- | -- |
| | \/ port 8080 to admin fe and be via /api \ | |
| route-- |                    |-- MongoDB |
| | \ prod 80 to production fe and be via /api / | |

## What we are using in OpenShift
1. Internal container registry for storing our three images
2. Configmap - we overwrite config values in both backend to connect to our DB and switch the role, because this is the same image.
3. MongoDB from the marketplace
4. The __route__ as native loadbalancer

## The installation sequence
- install [minishift](https://docs.okd.io/3.11/minishift/getting-started/installing.html)
- run minishift
- create the new project __myproject__
- to check if all working properly create first application __MongoDB__ with the standard parameters and __sampledb__ as database 
- to adjust the mongodb connect to pod 
oc exec mongodb-2-8m88d -i -t -- bash -il
login to mongo
mongo 127.0.0.1:27017/$MONGODB_DATABASE -u $MONGODB_USER -p $MONGODB_PASSWORD


- switch our docker to minishift reository 
<code><p>
$ minishift docker-env<p>
and execute all command in our shell. To verify -<p> 
$ docker ps
</code>
- [build all our images to openshift](https://docs.okd.io/3.11/minishift/using/docker-daemon.html). Or use save and load to push. <p>
We should yield the following in minishift repository

|REPOSITORY|TAG|IMAGE ID|CREATED|SIZE|
|----------|---|--------|-------|----|
|poc-cna-be    |      0.1       |          a49a4db8ffd8   |     30 hours ago   |     39.9MB|
|cna-admin-fe  |      0.1       |          345aa5c56e56   |     30 hours ago    |    24.1MB|
|cna-prod-fe   |      0.1       |          2a96f289fd7b   |     30 hours ago     |   24.1MB|

- login to oc as a developer
<code>oc login -u developer -p developer</code>
- create the configmaps for our backend - for admin and for prod
<code><p>
$ cd poc_cna_openshift/App/BackEnd<p>
$ oc create configmap admin-config --from-file=configs/cna-config.toml<p>
$ oc create configmap prod-config  --from-file=configs/cna-config.toml<p>
</code>
Amend in both files via web interface role, userpass what your specified during creating and in connection string replace the __localhost__ to mnongodb service name. Helpful command $ kubectl get services
Check the result $ oc describe configmap admin-config
-- deploy our applications
cd yaml
oc new-app --file=admin-backend-deploy.yaml
oc new-app --file=prod-backend-deploy.yaml
oc new-app --file=admin-frontend-deploy.yaml
oc new-app --file=prod-frontend-deploy.yaml

if something went wromg delete apps by the following command for axample
oc delete all -l app=poc-cna-be

to be continue...