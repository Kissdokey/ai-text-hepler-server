package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func Init() {
	cfg := weaviate.Config{
		Host:       "ai-text-helper-ptxsqxyy.weaviate.network", // Replace with your Weaviate endpoint
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: "fOcNXWGNczMEhaliDKEcYOeVz0gjxFCKOpff"}, // Replace with your Weaviate instance API key
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	isReady, err := client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", isReady)

	// Retrieve the data
	data, err := http.DefaultClient.Get("https://raw.githubusercontent.com/weaviate-tutorials/quickstart/main/data/jeopardy_tiny.json")
	if err != nil {
		panic(err)
	}
	defer data.Body.Close()

	// Decode the data
	var items []map[string]string
	if err := json.NewDecoder(data.Body).Decode(&items); err != nil {
		panic(err)
	}

	// convert items into a slice of models.Object
	objects := make([]*models.Object, len(items))
	for i := range items {
		objects[i] = &models.Object{
			Class: "Question",
			Properties: map[string]any{
				"category": items[i]["Category"],
				"question": items[i]["Question"],
				"answer":   items[i]["Answer"],
			},
		}
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			panic(res.Result.Errors.Error)
		}
	}

	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v111", result)
}
