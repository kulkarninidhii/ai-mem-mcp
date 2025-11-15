package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nidhi/ai-memory-mcp/go-server/memory"
)

func main() {
	store := memory.NewInMemStore()
	server := mcp.NewServer(&mcp.Implementation{Name: "memory-server", Version: "0.0.1"}, nil)

	type QueryInput struct {
		UserID string        `json:"user_id"`
		Query  string        `json:"query"`
		K      int           `json:"k"`
		Kinds  []memory.Kind `json:"kinds"`
	}
	type QueryOutput struct {
		Memories []*memory.Memory `json:"memories"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "query_memories",
		Description: "retrieve relevant memories for a user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input QueryInput) (*mcp.CallToolResult, QueryOutput, error) {
		if input.K == 0 {
			input.K = 8
		}
		mems, err := store.QueryRelevant(input.UserID, input.Query, input.K, input.Kinds)
		if err != nil {
			return nil, QueryOutput{}, nil
		}
		if mems == nil {
			mems = []*memory.Memory{}
		}
		return nil, QueryOutput{Memories: mems}, nil
	})

	type AddItem struct {
		Kind       memory.Kind `json:"kind"`
		Content    string      `json:"content"`
		Importance float64     `json:"importance"`
		Tags       []string    `json:"tags"`
	}
	type AddInput struct {
		UserID string    `json:"user_id"`
		Items  []AddItem `json:"items"`
	}
	type AddOutput struct {
		InsertedIDs []string `json:"inserted_ids"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_memories",
		Description: "store new memories for a user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AddInput) (*mcp.CallToolResult, AddOutput, error) {
		now := time.Now().UTC()
		ids := make([]string, 0, len(input.Items))

		for _, it := range input.Items {
			id := uuid.NewString()
			m := &memory.Memory{
				ID:         id,
				UserID:     input.UserID,
				Kind:       it.Kind,
				Content:    it.Content,
				Importance: it.Importance,
				Tags:       it.Tags,
				CreatedAt:  now,
				LastUsedAt: now,
			}
			if err := store.Insert(m); err != nil {
				return nil, AddOutput{}, nil
			}
			ids = append(ids, id)
		}

		return nil, AddOutput{InsertedIDs: ids}, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
