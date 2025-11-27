package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/juniorAkp/easyPay/pkg/types"
	"google.golang.org/genai"
)

func ExtractMessageDetails(ctx context.Context, userInput string) (*types.MessageDetails, error) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	if os.Getenv("GEMINI_API_KEY") == "" {
		fmt.Println("GEMINI_API_KEY environment variable not set")
		return nil, err
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GenAI client: %w", err)
	}

	systemInstruction := os.Getenv("PROMPT")

	prompt := fmt.Sprintf("%s\n\nUser message: %s", systemInstruction, userInput)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating content: %w", err)
	}

	responseText := result.Text()

	responseText = strings.TrimSpace(responseText)

	if idx := strings.Index(responseText, "{"); idx != -1 {
		responseText = responseText[idx:]
	}
	if idx := strings.LastIndex(responseText, "}"); idx != -1 {
		responseText = responseText[:idx+1]
	}

	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var details types.MessageDetails
	if err := json.Unmarshal([]byte(responseText), &details); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w (response: %s)", err, responseText)
	}

	return &details, nil
}
