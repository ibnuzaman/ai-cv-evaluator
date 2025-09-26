# Troubleshooting Guide - AI CV Evaluator

## Common Issues and Solutions

### 1. Gemini API Model Not Found Error

**Error:** `Publisher Model 'gemini-1.5-flash' was not found or your project does not have access to it`

**Solution:** 
- The application now uses `gemini-pro` which is widely available
- Ensure your Gemini API key is valid and has access to the Generative AI API
- Test your API key with: `./test-gemini`

### 2. Gemini API Key Issues

**Error:** `Gemini API key is required`

**Solution:**
1. Make sure you have a valid Google AI Studio API key
2. Add it to your `.env` file:
   ```
   GEMINI_API_KEY=your_api_key_here
   ```
3. Test the key with: `./test-gemini`

### 2.1. Gemini API Quota Exceeded

**Error:** `You exceeded your current quota, please check your plan and billing details`

**Solution:**
1. Check your Google AI Studio usage and billing
2. Wait for the quota to reset (usually 24 hours for free tier)
3. Consider upgrading to a paid plan for higher limits
4. The error message will show how long to wait (e.g., "Please retry in 53s")

### 3. ChromaDB Connection Issues

**Error:** `Warning: Failed to connect to ChromaDB`

**Solution:**
1. Start ChromaDB service: `docker-compose up -d chroma_server`
2. Check if ChromaDB is running: `curl http://localhost:8000/api/v2/heartbeat`
3. The application will continue without ChromaDB if it's not available

### 4. File Upload Issues

**Error:** `failed to read CV file` or `failed to read report file`

**Solution:**
1. Ensure the `uploads` directory exists and is writable
2. Check file formats - currently supports PDF and TXT files
3. Verify file paths are correct in the database

### 5. Database Connection Issues

**Error:** `Failed to connect to postgresql database`

**Solution:**
1. Start PostgreSQL: `docker-compose up -d db`
2. Check your `.env` database configuration
3. Ensure database exists and user has proper permissions

## Testing Commands

```bash
# Test full pipeline
./test-pipeline.sh

# Test Gemini API specifically
./test-gemini

# Check available Gemini models
./check-models.sh

# Test ChromaDB seeding
./seed

# Start all services
docker-compose up -d

# Check service status
docker-compose ps
```

## Model Alternatives

If `gemini-1.5-pro` doesn't work, you can try these alternatives in `internal/ai/gemini_client.go`:

- `gemini-1.0-pro` - Older but stable version  
- `gemini-pro` - Alternative naming (may work with some API keys)
- Check available models with: `curl -H "x-goog-api-key: $GEMINI_API_KEY" https://generativelanguage.googleapis.com/v1beta/models`

## API Rate Limits

Gemini API has rate limits:
- Free tier: Limited requests per minute
- Paid tier: Higher limits

Consider implementing exponential backoff for production use.
