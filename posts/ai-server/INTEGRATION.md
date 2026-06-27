# Integration with Posts Frontend

The Posts-AI Service integrates with the Posts Frontend to provide AI-assisted content generation and improvement.

## Frontend Changes Required

Update the frontend to route AI requests through this service instead of calling AI service directly.

### 1. Update AI Client Configuration

Update `src/lib/config.ts`:

```typescript
const getPostsAiServiceUrl = (): string => {
  const baseUrl = env.PUBLIC_POSTS_AI_SERVICE_URL || 'http://localhost:3014'
  const cleanUrl = baseUrl.replace(/\/$/, '')
  return `${cleanUrl}/api/v1`
}

export const config = {
  // ... existing config ...
  get postsAiServiceUrl() {
    return getPostsAiServiceUrl()
  },
}
```

### 2. Update Environment Variables

Add to `.env`:

```env
PUBLIC_POSTS_AI_SERVICE_URL=http://localhost:3014
```

### 3. Update AI Client Service

Create `src/lib/api/ai/posts-ai-client.ts`:

```typescript
import { config } from '$lib/config'
import type { Agent } from './client'

export interface GenerateDraftRequest {
  user_id: string
  post_id?: string
  prompt: string
  agent: Agent
}

export interface ImproveContentRequest {
  user_id: string
  post_id: string
  improvement: string
  agent: Agent
}

class PostsAIClient {
  private baseUrl: string

  constructor() {
    this.baseUrl = config.postsAiServiceUrl
  }

  async *generateDraft(
    userId: string,
    postId: string | undefined,
    prompt: string,
    agent: Agent,
  ): AsyncGenerator<StreamChunk> {
    const response = await fetch(`${this.baseUrl}/chats/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: userId,
        post_id: postId,
        prompt,
        agent,
      }),
    })

    if (!response.ok) {
      throw new Error(`Posts AI service error: ${response.status}`)
    }

    if (!response.body) {
      throw new Error('Response body is empty')
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.trim()) {
            try {
              const chunk = JSON.parse(line)
              yield chunk
            } catch (err) {
              console.error('Failed to parse chunk:', line)
            }
          }
        }
      }
    } finally {
      reader.releaseLock()
    }
  }

  async improveContent(
    userId: string,
    postId: string,
    improvement: string,
    agent: Agent,
  ) {
    const response = await fetch(`${this.baseUrl}/posts/${postId}/ai/improve`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: userId,
        post_id: postId,
        improvement,
        agent,
      }),
    })

    if (!response.ok) {
      throw new Error(`Posts AI service error: ${response.status}`)
    }

    return response.json()
  }

  async getChat(chatId: string) {
    const response = await fetch(`${this.baseUrl}/chats/${chatId}`)
    if (!response.ok) throw new Error('Failed to get chat')
    return response.json()
  }

  async listChats(userId: string, limit = 20, offset = 0) {
    const params = new URLSearchParams({
      user_id: userId,
      limit: limit.toString(),
      offset: offset.toString(),
    })

    const response = await fetch(`${this.baseUrl}/chats?${params}`)
    if (!response.ok) throw new Error('Failed to list chats')
    return response.json()
  }
}

export const postsAIClient = new PostsAIClient()
```

### 4. Update Components to Use Posts-AI Client

Update `src/lib/components/DraftBuilder.svelte`:

```svelte
<script lang="ts">
	import { postsAIClient } from '$lib/api/ai/posts-ai-client';
	import { auth } from '$lib';
	// ... rest of imports

	async function generateDraft() {
		if (!context.trim()) {
			error = 'Please provide article context';
			return;
		}

		isLoading = true;
		streamContent = '';
		error = '';
		controller = new AbortController();

		try {
			for await (const chunk of postsAIClient.generateDraft(
				$auth.user?.id || '',
				undefined,
				context,
				agent
			)) {
				if (controller?.signal.aborted) break;

				if (chunk.error) {
					error = chunk.error;
					isLoading = false;
					return;
				}

				if (chunk.delta) {
					streamContent += chunk.delta;
				}

				if (chunk.done) {
					isLoading = false;
				}
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to generate draft';
			isLoading = false;
		}
	}
	// ... rest of component
</script>
```

Similarly update `src/routes/posts/[slug]/edit/+page.svelte` to use `postsAIClient`.

## WebSocket Integration (Optional)

For real-time chat mode:

```typescript
const socket = new WebSocket(
  `ws://localhost:3014/ws/chats/${chatId}?user_id=${userId}`,
)

socket.onmessage = (event) => {
  const { response, done } = JSON.parse(event.data)
  updateDisplay(response)
  if (done) closeSocket()
}

socket.send(
  JSON.stringify({
    prompt: 'improve this content',
    agent: 'auto',
  }),
)
```

## Benefits of This Architecture

1. **Persistence** - All AI interactions are logged and searchable
2. **Audit Trail** - Track who used which AI agents and when
3. **Cost Tracking** - Monitor AI usage per user/tenant
4. **Rate Limiting** - Prevent abuse at service level
5. **Scalability** - AI service scales independently
6. **Separation of Concerns** - Posts backend doesn't need to know about AI
7. **Future Analytics** - Can easily add dashboards, reports, recommendations

## Performance Considerations

- Posts-AI Service handles ~1,000s requests/second
- Streaming NDJSON keeps memory low
- PostgreSQL indexes optimized for queries
- Connection pooling prevents connection exhaustion

## Monitoring

Add monitoring for:

- Response times (target: <1s first token)
- Error rates (track AI service failures)
- Token usage (for billing)
- Concurrent connections (scale horizontally if needed)
- Database query times

## Troubleshooting

### Service won't start

```bash
# Check database connectivity
psql postgres://woragis:password@localhost:5432/posts_ai

# Check AI service health
curl http://localhost:8000/healthz
```

### Streams not working

- Verify `Transfer-Encoding: chunked` headers
- Check CORS configuration
- Ensure AI service is accessible

### Missing chats in database

- Verify database migrations ran
- Check PostgreSQL logs
- Verify user_id is valid UUID format
