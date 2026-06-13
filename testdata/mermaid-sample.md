# Mermaid Diagram Sample

In the app, Mermaid blocks show as a labeled source block. Export to HTML
(or Print via Browser / PDF) to see them rendered as diagrams.

## Flowchart

```mermaid
graph TD;
    A[Start] --> B{Is it working?};
    B -- Yes --> C[Ship it];
    B -- No --> D[Fix it];
    D --> B;
```

## Sequence diagram

```mermaid
sequenceDiagram
    User->>MarkView: Open document
    MarkView->>Renderer: Parse markdown
    Renderer-->>User: Rendered view
```

## Regular code (not a diagram)

```go
fmt.Println("this stays a normal code block")
```
