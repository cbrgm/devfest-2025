package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tea represents a tea item in our shop
type Tea struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	InStock     bool    `json:"in_stock"`
}

// Tea inventory
var teaInventory = []Tea{
	{ID: 1, Name: "Earl Grey", Type: "Black", Price: 12.99, Description: "Classic bergamot-scented black tea", InStock: true},
	{ID: 2, Name: "Dragon Well", Type: "Green", Price: 18.50, Description: "Delicate Chinese green tea with a nutty flavor", InStock: true},
	{ID: 3, Name: "Chamomile Dream", Type: "Herbal", Price: 14.25, Description: "Soothing caffeine-free chamomile blend", InStock: true},
	{ID: 4, Name: "Royal Pu-erh", Type: "Pu-erh", Price: 25.00, Description: "Aged fermented tea with deep, earthy flavor", InStock: false},
	{ID: 5, Name: "Jasmine Phoenix Pearls", Type: "Green", Price: 32.00, Description: "Hand-rolled green tea scented with jasmine flowers", InStock: true},
}

// Tool parameter structs with jsonschema tags for automatic schema inference
// Fields without omitempty in the json tag are required
type GetTeaParams struct {
	ID int `json:"id" jsonschema:"The ID of the tea to retrieve"`
}

type SearchTeasParams struct {
	Query string `json:"query" jsonschema:"Search term for tea name, type, or description"`
}

type CheckStockParams struct {
	ID int `json:"id" jsonschema:"The ID of the tea to check stock for"`
}

// List teas handler - takes empty struct for no parameters
func ListTeas(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	teaJSON, _ := json.MarshalIndent(teaInventory, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Available teas in our shop:\n%s", string(teaJSON))}},
	}, nil, nil
}

// Get tea handler
func GetTea(ctx context.Context, req *mcp.CallToolRequest, args GetTeaParams) (*mcp.CallToolResult, any, error) {
	teaID := args.ID

	for _, tea := range teaInventory {
		if tea.ID == teaID {
			teaJSON, _ := json.MarshalIndent(tea, "", "  ")
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Tea details:\n%s", string(teaJSON))}},
			}, nil, nil
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Tea with ID %d not found", teaID)}},
		IsError: true,
	}, nil, nil
}

// Search teas handler
func SearchTeas(ctx context.Context, req *mcp.CallToolRequest, args SearchTeasParams) (*mcp.CallToolResult, any, error) {
	query := strings.ToLower(args.Query)
	var matchingTeas []Tea

	for _, tea := range teaInventory {
		if strings.Contains(strings.ToLower(tea.Name), query) ||
			strings.Contains(strings.ToLower(tea.Type), query) ||
			strings.Contains(strings.ToLower(tea.Description), query) {
			matchingTeas = append(matchingTeas, tea)
		}
	}

	if len(matchingTeas) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No teas found matching '%s'", args.Query)}},
		}, nil, nil
	}

	resultJSON, _ := json.MarshalIndent(matchingTeas, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Found %d tea(s) matching '%s':\n%s", len(matchingTeas), args.Query, string(resultJSON))}},
	}, nil, nil
}

// Check stock handler
func CheckStock(ctx context.Context, req *mcp.CallToolRequest, args CheckStockParams) (*mcp.CallToolResult, any, error) {
	teaID := args.ID

	for _, tea := range teaInventory {
		if tea.ID == teaID {
			status := "In Stock"
			if !tea.InStock {
				status = "Out of Stock"
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%s (%s) - %s", tea.Name, tea.Type, status)}},
			}, nil, nil
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Tea with ID %d not found", teaID)}},
		IsError: true,
	}, nil, nil
}

// Resource handler for tea catalog
func GetTeaCatalog(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	catalogJSON, _ := json.MarshalIndent(teaInventory, "", "  ")
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     string(catalogJSON),
		}},
	}, nil
}

func main() {
	// Create server with implementation details and options
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "teashop-mcp-server",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		Instructions: "Tea Shop MCP Server - Browse our tea collection, get details, and check stock availability",
	})

	// Register tools with the server using the generic AddTool function
	// Schema is automatically inferred from the parameter types and jsonschema tags
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_teas",
		Description: "List all available teas in the shop",
	}, ListTeas)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_tea",
		Description: "Get detailed information about a specific tea by ID",
	}, GetTea)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_teas",
		Description: "Search for teas by name, type, or description",
	}, SearchTeas)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_stock",
		Description: "Check if a specific tea is in stock",
	}, CheckStock)

	// Register resource with the server
	server.AddResource(&mcp.Resource{
		URI:         "teashop://catalog",
		Name:        "Tea Catalog",
		Description: "Complete catalog of all teas available in the shop",
		MIMEType:    "application/json",
	}, GetTeaCatalog)

	log.Printf("Starting Tea Shop MCP Server v1.0.0")
	log.Printf("Registered 4 tools and 1 resource")

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
