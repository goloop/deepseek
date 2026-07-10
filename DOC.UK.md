# deepseek - довідник

Повний довідник пакета `deepseek`: клієнт, спільна модель `goloop/ai`, chat
completions (інтерфейс і нативний), стрімінг, reasoning-контент і моделі.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Створення клієнта](#створення-клієнта)
- [Generate і Stream](#generate-і-stream)
- [Нативні chat completions](#нативні-chat-completions)
- [Reasoning-моделі](#reasoning-моделі)
- [Інструменти й system-промпти](#інструменти-й-system-промпти)
- [Моделі](#моделі)
- [Опції та помилки](#опції-та-помилки)

## Ментальна модель

`deepseek.Client` реалізує `ai.Client` - провайдер-незалежний контракт із
`github.com/goloop/ai`. Спільні `Generate` і `Stream` покривають спільну основу
(чат із інструментами й стрімінгом), тож код проти інтерфейсу працює з будь-яким
провайдером.

Специфіка провайдера - у нативних методах: повний `ChatCompletion` і перелік
моделей. Формат обміну - сумісний із chat completions.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/deepseek"
)
```

## Створення клієнта

```go
c := deepseek.New(os.Getenv("DEEPSEEK_API_KEY"))

c = deepseek.New(apiKey, deepseek.WithTimeout(30*time.Second))
```

Base URL за замовчуванням `https://api.deepseek.com`. Наведіть `WithBaseURL` на
будь-який сумісний ендпоінт, щоб перевикористати клієнт.

## Generate і Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    deepseek.ModelChat,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` повертає `iter.Seq2[ai.Chunk, error]`: текстові дельти чанками з `Text`,
завершений виклик інструмента - чанком із `ToolCall`, фінальний чанк - `Done` і
`Usage`.

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		return err
	}
	fmt.Print(chunk.Text)
}
```

## Нативні chat completions

Для опцій, специфічних для провайдера, будуйте `ChatRequest` і викликайте
`ChatCompletion` чи `ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &deepseek.ChatRequest{
	Model:          deepseek.ModelChat,
	Messages:       []deepseek.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

Доступні `Tools`, `ToolChoice`, `Temperature`, `TopP`, `MaxTokens`, `Stop`, `N`,
`Seed`, `ResponseFormat`, `User`.

## Reasoning-моделі

`deepseek-reasoner` повертає ланцюг міркувань у `ReasoningContent`, окремо від
фінальної відповіді в `Content`:

```go
resp, _ := c.ChatCompletion(ctx, &deepseek.ChatRequest{
	Model:    deepseek.ModelReasoner,
	Messages: []deepseek.ChatMessage{{Role: "user", Content: "9.11 or 9.8?"}},
})
resp.Choices[0].Message.ReasoningContent // міркування
resp.Choices[0].Message.Content          // відповідь
```

`Generate` і `Stream` повертають лише текст фінальної відповіді. Нативний
`ChatStream` стрімить і міркування: кожна дельта `ChatStreamChunk` несе
`ReasoningContent` поряд із `Content`.

## Інструменти й system-промпти

Інструменти й system-промпти використовують спільні типи `ai`: `ai.Tool`,
`ai.ToolResult` і повідомлення `RoleSystem` або поле `System`. Результати
інструментів надсилаються назад повідомленнями `RoleTool`, де `ai.ToolResult.ID`
збігається з `ai.ToolUse.ID`.

## Моделі

```go
models, err := c.Models(ctx)
m, err := c.GetModel(ctx, deepseek.ModelChat)
```

## Опції та помилки

Опції: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`.

Невдала відповідь стає `*ai.APIError` зі `Status`, `Type`, `Code`, `Message` і
сирим тілом:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// backoff
}
```

Запити без моделі чи повідомлень падають до мережі з `ai.ErrNoModel` або
`ai.ErrNoMessages`.
