# Google

Relay can make use of Google PubSub and Google Datastore.
You must have a project with billing enabled, however,
relay's usage of PubSub and Datastore fits well within the
[Free Tier](https://cloud.google.com/free/docs/gcp-free-tier)
range unless you have a very busy deployment.

## Setup

Configure the project, datastore, pubsub, service account

### Project

Create the project and enable billing

***create project***
```
PROJECT_NAME=somename

gcloud projects create $PROJECT_NAME
gcloud config set project $PROJECT_NAME
```

***enable billing***
```
gcloud beta billing accounts list

BILLING=<your billing accoun idt>
```
```
gcloud beta billing projects link $PROJECT_NAME --billing-account $BILLING
```

### Datastore

Datastore is used for storing authentication tokens (UUIDs)
for each account.

Enable the datastore api, initialize cloud datastore.

1. In the web UI, [enable datastore mode for cloud datastore](https://cloud.google.com/firestore/docs/firestore-or-datastore?authuser=1&_ga=2.212492372.-1109796302.1520098783)
2. enable api
```
gcloud services enable datastore.googleapis.com
```
3. create your first account using the web ui
  - click "create entity"
  - kind: `account`
  - add property
    - name: APIKey
    - value: <use a random UUID>
  - add property
    - name: Name
    - value: <the name of your new account>

### PubSub

Enable the API, create a topic and subscription

```
gcloud services enable pubsub.googleapis.com
```
```
gcloud pubsub topics create <topic name>
gcloud pubsub subscriptions create --topic=<topic name> <subscription name>
```

### Service Account

A [service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-gcloud)
should be used for Relay's interaction with Google APIs.

create the service account and json key
```
PROJECT_ID=<your project ID>

gcloud iam service-accounts create [SA-NAME] \
    --description "[SA-DESCRIPTION]" \
    --display-name "[SA-DISPLAY-NAME]"

gcloud iam service-accounts keys create [OUTPUT-JSON-FILE-NAME] \
    --iam-account [SA-NAME]@${PROJECT_ID}.iam.gserviceaccount.com
```
setup IAM roles
```
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:[SA-NAME]@${PROJECT_ID}.iam.gserviceaccount.com \
    --role roles/datastore.user

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:[SA-NAME]@${PROJECT_ID}.iam.gserviceaccount.com \
    --role roles/pubsub.publisher

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:[SA-NAME]@${PROJECT_ID}.iam.gserviceaccount.com \
    --role roles/pubsub.subscriber
```
