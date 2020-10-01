# This repository contains the source files for POC - new application in OpenShift cluster.
## Components:
1. BackEnd golang application. We will use it on a two roles - as a client only and as a Server only. Purpose - emulate production and admin wing.
2. Database for storing our test data MongoDB - OpenShift application launched from the OS Marketplace. Nothing special - create new and delate existed dicuments.
3. FrontEnd - versy simple html pages. One as for admin page - able to create new or delete cna user. Second one - production page to viewonly model.
4. Last part - the folder with command for OpenShift to maniulate with the cluster (oc commands and yaml files)

## Architecture:
Inside one project/namespace we have 3 internal services/applications and the route, who exosing our solution and managing level 7 trafic.

| routes to | service name | db | 
| --                 |         -- |         -- | 
| admin.cna.com/     |  admin fe  |            |
| admin.cna.com/api  |  admin be  | -- MongoDB |
|                    |            |   same     |
| www.cna.com/api    |  prod be   | -- MongoDB |
| www.cna.com/       |  prod fe   |            |

## What we are using in OpenShift
1. Internal container registry for storing our three images
2. Configmap - we overwrite config values in both backend to connect to our DB and switch the role, because this is the same image.
3. MongoDB from the marketplace
4. The __route__ as native loadbalancer

## The installation sequence
- install [minishift](https://docs.okd.io/3.11/minishift/getting-started/installing.html)
- run minishift
- create the new project __myproject__
- to check if all working properly create first application __MongoDB__ with the following: __monogodb__ dbuser/dbpass, admin pass: __admin__ and __sampledb__ as database 
Connect to pod via terminal and execute the following:

<code>
mongo -u admin -p admin admin<p>
db.createUser({<p>
  user: "mongodb2",<p>
  pwd: "mongodb2",<p>
  roles: [<p>
    { role: "userAdmin", db: "sampledb" },<p>
    { role: "dbAdmin",   db: "sampledb" },<p>
    { role: "readWrite", db: "sampledb" }<p>
  ]<p>
});<p>
</code>
This is all with Mongo.
- To continue the deploy switch our docker to minishift reository 

<code><p>
$ minishift docker-env<p>
and execute all command in our shell. To verify -<p> 
$ docker ps
</code>
- [build all our images to openshift](https://docs.okd.io/3.11/minishift/using/docker-daemon.html). Or use save and load to push - I'm using for this two terminal window. <p>
We should yield the following in minishift repository

|REPOSITORY|TAG|IMAGE ID|CREATED|SIZE|
|----------|---|--------|-------|----|
|poc-cna-be    |      0.1       |          a49a4db8ffd8   |     30 hours ago   |     39.9MB|
|cna-admin-fe  |      0.1       |          345aa5c56e56   |     30 hours ago    |    24.1MB|
|cna-prod-fe   |      0.1       |          2a96f289fd7b   |     30 hours ago     |   24.1MB|
We are ready.
- login to oc as a developer

<code>oc login -u developer -p developer</code>
- create the configmaps for our backend - for admin and for prod
<code><p>
$ cd poc_cna_openshift/App/BackEnd<p>
$ oc create configmap admin-config --from-file=configs/cna-config.toml<p>
$ oc create configmap prod-config  --from-file=configs/cna-config.toml<p>
</code>

Amend in both files via web interface role, userpass what your specified during creating and in connection string replace the __localhost__ to mnongodb service name. Helpful command $ kubectl get services<p>
Check the result <p>

<code>
$ oc describe configmap admin-config<p>
</code>
- deploy our applications<p>
<code>
cd yaml<p>
oc create -f admin-backend-deploy.yaml<p>
oc create -f prod-backend-deploy.yaml<p>
oc create -f admin-frontend-deploy.yaml<p>
oc create -f prod-frontend-deploy.yaml<p>
</code>
- create our routes<p>
<code>
oc create -f routes-all.yaml<p>
</code>
- Cluster loadbalancer looking for name scpecific in routes - www.cna.com and admin.cna.com. Thus, we have to substitute this domain names to cluster IP.<p>
Edit /etc/hosts and add:<p>
<code>
#minishift<p>
192.168.99.103  admin.cna.com<p>
192.168.99.103  www.cna.com<p>
</code>

Finish.

Lets check:<p>
(http://admin.cna.com) - should display data grid<p>
(http://admin.cna.com/api) - Display that we are Admin.<p>
(http://admin.cna.com/api/get) - json with all documents data<p>
the same for Prod<p>
(http://www.cna.com) - datagrid in readonly mode<p>
(http://www.cna.com/api) - Display that we are Client.<p>
(http://www.cna.com/api/get) - json with all documents data<p>
