package kjarni

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// C struct layouts

type ffiClassifierConfig struct {
	Device    int32
	_         int32 // padding
	CacheDir  uintptr
	ModelName uintptr
	ModelPath uintptr
	Labels    uintptr
	NumLabels uintptr
	MultiLabel int32
	Quiet     int32
}

type ffiClassResult struct {
	Label uintptr
	Score float32
	_     [4]byte // padding to 16 bytes
}

type ffiClassResults struct {
	Results uintptr
	Len     uintptr
}

type ffiEmbedderConfig struct {
	Device    int32
	_         int32 // padding
	CacheDir  uintptr
	ModelName uintptr
	ModelPath uintptr
	Normalize int32
	Quiet     int32
}

type ffiFloatArray struct {
	Data uintptr
	Len  uintptr
}

type ffiFloat2DArray struct {
	Data uintptr
	Rows uintptr
	Cols uintptr
}

type ffiRerankerConfig struct {
	Device    int32
	_         int32 // padding
	CacheDir  uintptr
	ModelName uintptr
	ModelPath uintptr
	Quiet     int32
	_2        int32 // padding
}

type ffiRerankResult struct {
	Index uintptr
	Score float32
	_     [4]byte // padding
}

type ffiRerankResults struct {
	Results uintptr
	Len     uintptr
}

var (
	ffiOnce sync.Once

	// Error handling
	_lastErrorMessage func() uintptr
	_clearError       func()

	// Classifier
	_classifierNewSym  uintptr
	_classifierFree    func(handle uintptr)
	_classifierClassifySym uintptr
	_classifierNumLabels func(handle uintptr) uintptr
	_classResultsFree  func(results uintptr, len uintptr)

	// Embedder
	_embedderNewSym    uintptr
	_embedderFree      func(handle uintptr)
	_embedderEncodeSym uintptr
	_embedderEncodeBatchSym uintptr
	_embedderSimilaritySym uintptr
	_embedderDim       func(handle uintptr) uintptr
	_floatArrayFree    func(data uintptr, len uintptr)
	_float2DArrayFree  func(data uintptr, rows uintptr, cols uintptr)

	// Reranker
	_rerankerNewSym    uintptr
	_rerankerFree      func(handle uintptr)
	_rerankerScoreSym  uintptr
	_rerankerRerankSym uintptr
	_rerankerRerankTopKSym uintptr
	_rerankResultsFree func(results uintptr, len uintptr)

	// Indexer
	_indexerNewSym    uintptr
	_indexerFree      func(handle uintptr)
	_indexerCreateSym uintptr

	// Searcher
	_searcherNewSym              uintptr
	_searcherFree                func(handle uintptr)
	_searcherSearchWithOptionsSym uintptr

	// Cosine similarity
	_cosineSimilarity func(a unsafe.Pointer, b unsafe.Pointer, len uintptr) float32

)

func initFFI() error {
	handle, err := loadLibrary()
	if err != nil {
		return err
	}

	// Error handling
	sym, err := findSymbol(handle, "kjarni_last_error_message")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_lastErrorMessage, sym)

	sym, err = findSymbol(handle, "kjarni_clear_error")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_clearError, sym)

	// Classifier
	_classifierNewSym, err = findSymbol(handle, "kjarni_classifier_new")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_classifier_free")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_classifierFree, sym)

	_classifierClassifySym, err = findSymbol(handle, "kjarni_classifier_classify")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_classifier_num_labels")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_classifierNumLabels, sym)

	_classResultsFreeSymGlobal, err = findSymbol(handle, "kjarni_class_results_free")
	if err != nil {
		return err
	}

	// Embedder
	_embedderNewSym, err = findSymbol(handle, "kjarni_embedder_new")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_embedder_free")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_embedderFree, sym)

	_embedderEncodeSym, err = findSymbol(handle, "kjarni_embedder_encode")
	if err != nil {
		return err
	}

	_embedderEncodeBatchSym, err = findSymbol(handle, "kjarni_embedder_encode_batch")
	if err != nil {
		return err
	}

	_embedderSimilaritySym, err = findSymbol(handle, "kjarni_embedder_similarity")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_embedder_dim")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_embedderDim, sym)

	_floatArrayFreeSym, err = findSymbol(handle, "kjarni_float_array_free")
	if err != nil {
		return err
	}

	_float2DArrayFreeSym, err = findSymbol(handle, "kjarni_float_2d_array_free")
	if err != nil {
		return err
	}

	// Reranker
	_rerankerNewSym, err = findSymbol(handle, "kjarni_reranker_new")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_reranker_free")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_rerankerFree, sym)

	_rerankerScoreSym, err = findSymbol(handle, "kjarni_reranker_score")
	if err != nil {
		return err
	}

	_rerankerRerankSym, err = findSymbol(handle, "kjarni_reranker_rerank")
	if err != nil {
		return err
	}

	_rerankerRerankTopKSym, err = findSymbol(handle, "kjarni_reranker_rerank_top_k")
	if err != nil {
		return err
	}

	_rerankResultsFreeSym, err = findSymbol(handle, "kjarni_rerank_results_free")
	if err != nil {
		return err
	}

	// Indexer
	_indexerNewSym, err = findSymbol(handle, "kjarni_indexer_new")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_indexer_free")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_indexerFree, sym)

	_indexerCreateSym, err = findSymbol(handle, "kjarni_indexer_create")
	if err != nil {
		return err
	}

	// Searcher
	_searcherNewSym, err = findSymbol(handle, "kjarni_searcher_new")
	if err != nil {
		return err
	}

	sym, err = findSymbol(handle, "kjarni_searcher_free")
	if err != nil {
		return err
	}
	purego.RegisterFunc(&_searcherFree, sym)

	_searcherSearchWithOptionsSym, err = findSymbol(handle, "kjarni_searcher_search_with_options")
	if err != nil {
		return err
	}

	_searchResultsFreeSym, err = findSymbol(handle, "kjarni_search_results_free")
	if err != nil {
		return err
	}

	return nil
}
// Symbols stored for SyscallN usage
var (
	_classResultsFreeSymGlobal uintptr
	_floatArrayFreeSym         uintptr
	_float2DArrayFreeSym       uintptr
	_rerankResultsFreeSym      uintptr
	_searchResultsFreeSym      uintptr
)

// convert Go string to null-terminated C string, returns pointer and cleanup func
func cString(s string) (uintptr, func()) {
	b := append([]byte(s), 0)
	ptr := unsafe.Pointer(&b[0])
	// prevent GC from collecting the slice while we use the pointer
	return uintptr(ptr), func() {
		_ = b
	}
}

// get last error message from FFI
func lastError(code int32) error {
	ptr := _lastErrorMessage()
	if ptr == 0 {
		return &KjarniError{
			Code:    ErrorCode(code),
			Message: "unknown error",
		}
	}
	msg := goString(ptr)
	return &KjarniError{
		Code:    ErrorCode(code),
		Message: msg,
	}
}

// read null-terminated C string from pointer
func goString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	var length int
	for {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(length)))
		if b == 0 {
			break
		}
		length++
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
	}
	return string(bytes)
}