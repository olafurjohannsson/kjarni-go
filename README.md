# kjarni-go

Text classification, embeddings, semantic search, and reranking for Go. No Python, no containers, no ONNX. One `go get` and you're running inference.

```
go get github.com/olafurjohannsson/kjarni-go
```

# Kjarni inference engine

This go library uses the custom built [Kjarni](https://github.com/olafurjohannsson/kjarni) engine to do inference

Models download automatically on first use and are cached locally.

## Classify

```go
c, _ := kjarni.NewClassifier("roberta-sentiment", kjarni.WithQuiet(true))
defer c.Close()

result, _ := c.Classify("I love this product!")
fmt.Println(result)
// positive (98.5%)
```

Available models: `distilbert-sentiment`, `roberta-sentiment`, `bert-sentiment-multilingual`, `distilroberta-emotion`, `roberta-emotions`, `toxic-bert`

## Embeddings

```go
e, _ := kjarni.NewEmbedder("minilm-l6-v2", kjarni.WithQuiet(true))
defer e.Close()

sim, _ := e.Similarity("doctor", "physician")
fmt.Printf("%.4f\n", sim)
// 0.8598

vec, _ := e.Encode("hello world")
fmt.Println(len(vec))
// 384

vecs, _ := e.EncodeBatch([]string{"cat", "dog", "airplane"})
```

Available models: `minilm-l6-v2` (384d), `mpnet-base-v2` (768d), `distilbert-base` (768d)

## Search

Index a directory and search using keyword (BM25), semantic (vector), or hybrid (both combined).

```go
// index
idx, _ := kjarni.NewIndexer("minilm-l6-v2", kjarni.WithQuiet(true))
idx.Create("/path/to/index", []string{"/path/to/docs"})
idx.Close()

// search
s, _ := kjarni.NewSearcher("minilm-l6-v2", "", kjarni.WithQuiet(true))
defer s.Close()

results, _ := s.Search("/path/to/index", "how do returns work?", kjarni.Hybrid)
for _, r := range results {
    fmt.Printf("%.4f: %s\n", r.Score, r.Text)
}
```

To enable cross-encoder reranking, pass a reranker model when creating the searcher:

```go
s, _ := kjarni.NewSearcher("minilm-l6-v2", "minilm-l6-v2-cross-encoder", kjarni.WithQuiet(true))
```

## Rerank

Score and sort documents by relevance to a query using a cross-encoder.

```go
r, _ := kjarni.NewReranker(kjarni.WithQuiet(true))
defer r.Close()

docs := []string{
    "The weather is nice today.",
    "Machine learning is a branch of AI.",
    "Deep learning uses neural networks.",
}

ranked, _ := r.Rerank("What is machine learning?", docs)
for _, result := range ranked {
    fmt.Printf("[%d] %.2f  %s\n", result.Index, result.Score, result.Document)
}

// or get only the top k
top, _ := r.RerankTopK("machine learning", docs, 1)
```

## How it works

This package embeds a Rust inference engine as a shared library (`.so` on Linux, `.dll` on Windows). The library is extracted to a temp directory at runtime and loaded via [purego](https://github.com/ebitengine/purego) — no cgo required.

The same engine powers the [C# NuGet package](https://www.nuget.org/packages/Kjarni), the CLI, and the WASM build.

## Platform support

| OS | Arch | Status |
|----|------|--------|
| Linux | amd64 | supported |
| Windows | amd64 | supported |

## License

MIT
