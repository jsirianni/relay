package gcpdatastore

import (
    "context"

    "cloud.google.com/go/datastore"
    "google.golang.org/api/iterator"

    "github.com/pkg/errors"
    "github.com/google/uuid"
)

type Datastore struct {
    ProjectID string
    Client    *datastore.Client

    ctx       context.Context
}

type Account struct {
    Name string
    APIKey  string
}

const DatastoreKind = "Account"

func New(projectID string) (Datastore, error) {
    if projectID == "" {
        return Datastore{}, errors.New("project id is not set, cannot init Datastore type")
    }

    var err error
    d := Datastore{
        ProjectID: projectID,
        ctx: context.Background(),
    }

    d.Client, err = datastore.NewClient(d.ctx, d.ProjectID)
    return d, err
}

func (d Datastore) ValidAPIKey(apiKey uuid.UUID) (bool, error) {
    q := datastore.NewQuery(DatastoreKind).Filter("APIKey =", apiKey.String())
    it := d.Client.Run(d.ctx, q)
    for {
        var a Account
        _, err := it.Next(&a)
        if err == iterator.Done {
            return false, nil
        }
        if err != nil {
            return false, err
        }

        if a.APIKey == apiKey.String() {
            return true, nil
        }
    }

    return false, nil
}

func (d Datastore) Type() string {
    return "cloud-datastore"
}
