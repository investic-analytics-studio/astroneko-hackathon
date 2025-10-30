package history

import (
	"encoding/json"
	"regexp"
	"strings"
)

type CardData struct {
	Card    string `json:"card"`
	Meaning string `json:"meaning"`
}

func ExtractJSONFromMessage(message string) (cleanedMessage string, card string, meaning string) {
	cleanedMessage = message
	card = ""
	meaning = ""

	jsonBlockRegex := regexp.MustCompile("(?s)```json\\s*(.+?)\\s*```")
	matches := jsonBlockRegex.FindStringSubmatch(message)

	if len(matches) > 1 {
		jsonContent := strings.TrimSpace(matches[1])

		var cardData CardData
		if err := json.Unmarshal([]byte(jsonContent), &cardData); err == nil {
			card = cardData.Card
			meaning = cardData.Meaning

			cleanedMessage = jsonBlockRegex.ReplaceAllString(message, "")
			cleanedMessage = strings.TrimSpace(cleanedMessage)
		} else {
			// Try to fix malformed JSON by adding braces if missing
			fixedJSON := tryFixMalformedJSON(jsonContent)
			if fixedJSON != "" {
				if err := json.Unmarshal([]byte(fixedJSON), &cardData); err == nil {
					card = cardData.Card
					meaning = cardData.Meaning

					cleanedMessage = jsonBlockRegex.ReplaceAllString(message, "")
					cleanedMessage = strings.TrimSpace(cleanedMessage)
				}
			}
		}
	}

	return cleanedMessage, card, meaning
}

// tryFixMalformedJSON attempts to fix common JSON formatting issues
func tryFixMalformedJSON(jsonContent string) string {
	// Remove any leading/trailing non-JSON characters
	trimmed := strings.TrimSpace(jsonContent)

	// If it doesn't start with {, try to wrap it
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		// Check if it looks like JSON key-value pairs without braces
		if strings.Contains(trimmed, ":") && (strings.Contains(trimmed, "\"card\"") || strings.Contains(trimmed, "\"meaning\"")) {
			return "{" + trimmed + "}"
		}
	}

	return ""
}
