package writer

import (
	wmetrics "github.com/YandexClassifieds/jaeger-ydb-store/storage/spanstore/writer/metrics"
	"github.com/uber/jaeger-lib/metrics"
	"sync"
)

type batchWriterMetrics struct {
	traces       *wmetrics.WriteMetrics
	spansDropped metrics.Counter
}

func newBatchWriterMetrics(factory metrics.Factory) batchWriterMetrics {
	return batchWriterMetrics{
		traces:       wmetrics.NewWriteMetrics(factory, "traces"),
		spansDropped: factory.Counter(metrics.Options{Name: "spans_dropped"}),
	}
}

type spanMetricsKey struct {
	svc string
	op  string
}

type invalidSpanMetrics struct {
	mf metrics.Factory
	m  map[spanMetricsKey]metrics.Counter
	mu sync.Mutex
}

func newInvalidSpanMetrics(mf metrics.Factory) *invalidSpanMetrics {
	return &invalidSpanMetrics{
		mf: mf,
		m:  make(map[spanMetricsKey]metrics.Counter, 0),
	}
}

func (m *invalidSpanMetrics) Inc(svc, op string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := spanMetricsKey{svc: svc, op: op}
	if _, exists := m.m[k]; !exists {
		m.m[k] = m.mf.Counter(metrics.Options{Name: "invalid_spans", Tags: map[string]string{"svc": svc, "op": op}})
	}
	m.m[k].Inc(1)
}
