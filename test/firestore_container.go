package test

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupFirestoreContainer(t *testing.T, projectId string) *firestore.Client {
	t.Helper()
	ctx := context.Background()

	firestoreContainer := runFirestoreContainer(t, ctx, projectId)
	client := createFirestoreClient(t, ctx, firestoreContainer)

	t.Cleanup(func() {
		cleanupFirestoreClient(t, client)
		terminateFirestoreContainer(t, firestoreContainer)
	})

	return client
}

func runFirestoreContainer(t *testing.T, ctx context.Context, projectId string) *gcloud.GCloudContainer {
	firestoreContainer, err := gcloud.RunFirestore(
		ctx,
		"gcr.io/google.com/cloudsdktool/cloud-sdk:513.0.0-emulators",
		gcloud.WithProjectID(projectId),
	)
	if err != nil {
		t.Fatalf("failed to run container: %v", err)
	}
	return firestoreContainer
}

func createFirestoreClient(t *testing.T, ctx context.Context, firestoreContainer *gcloud.GCloudContainer) *firestore.Client {
	conn := createGRPCConnection(t, firestoreContainer)
	projectID := firestoreContainer.Settings.ProjectID
	options := []option.ClientOption{option.WithGRPCConn(conn)}
	client, err := firestore.NewClient(ctx, projectID, options...)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func createGRPCConnection(t *testing.T, firestoreContainer *gcloud.GCloudContainer) *grpc.ClientConn {
	conn, err := grpc.NewClient(firestoreContainer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create grpc connection: %v", err)
	}
	return conn
}

func cleanupFirestoreClient(t *testing.T, client *firestore.Client) {
	t.Logf("firestore client cleanup")
	client.Close()
}

func terminateFirestoreContainer(t *testing.T, firestoreContainer *gcloud.GCloudContainer) {
	t.Logf("terminating container")
	if err := testcontainers.TerminateContainer(firestoreContainer); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}
