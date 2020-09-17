package v2api

import (
	"context"
	"log"

	"gitlab.com/abios/user-svc/pkg/database/queries"
	"gitlab.com/abios/user-svc/pkg/database/txhelper"
	"gitlab.com/abios/user-svc/pkg/structs"
	"gitlab.com/abios/user-svc/pkg/wraperr"
)

// updateApiClientResources updates the resources for the client and game pair
// specified by the passed API subscription struct; this means that the client
// will only be able to access resources specified by the pair's current
// package configuration.
func updateApiClientResources(ctx context.Context, subscription structs.ApiSubscription) error {
	// Set up transaction for updating the client's v3 resources.
	log.Println("transaction")
	tx, err := txhelper.Transaction()
	if err != nil {
		return wraperr.New(err).Errorln("Unable to start transaction for customer resource update")
	}

	// Instead of adding and removing the specific resources, of which the diff is
	// difficult to calculate, we instead remove all resources and add the
	// appropriate ones anew.
	log.Println("delete")
	if err := queries.DeleteApiClientResources(tx, ctx, subscription); err != nil {
		return wraperr.New(txhelper.Rollback(tx, err)).Errorln("Unable to remove client's existing resources")
	}

	// Figure out which resources the current API client subscription should have.
	log.Println("resolve")
	resources, err := queries.ResolveApiClientResources(tx, ctx, subscription)
	if err != nil {
		return wraperr.New(txhelper.Rollback(tx, err)).Errorln("Unable to resolve client's existing packages")
	}

	// Add the appropriate v3 resources to the (client + game) pair
	log.Println("add")
	if err := queries.AddApiClientResources(tx, ctx, subscription, resources); err != nil {
		return wraperr.New(txhelper.Rollback(tx, err)).Errorln("Unable to remove client's existing resources")
	}

	// Commit the transaction
	log.Println("commit")
	if err := txhelper.Commit(tx, nil); err != nil {
		return wraperr.New(err).Errorln("Unable to commit transaction for customer resource update")
	}

	return nil
}
