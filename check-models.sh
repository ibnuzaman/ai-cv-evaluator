#!/bin/bash

echo "Checking available Gemini API models..."

if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY environment variable not set"
    echo "Please set it in your .env file or export it"
    exit 1
fi

echo "🔍 Fetching available models from Gemini API..."

curl -s -H "x-goog-api-key: $GEMINI_API_KEY" \
     "https://generativelanguage.googleapis.com/v1beta/models" | \
     jq -r '.models[] | select(.supportedGenerationMethods[]? == "generateContent") | .name' 2>/dev/null || \
     curl -s -H "x-goog-api-key: $GEMINI_API_KEY" \
          "https://generativelanguage.googleapis.com/v1beta/models"

echo ""
echo "💡 Tip: Use one of the models that supports 'generateContent' method"
echo "Common working models:"
echo "  - models/gemini-1.5-pro"
echo "  - models/gemini-1.0-pro"
echo "  - models/gemini-pro"
