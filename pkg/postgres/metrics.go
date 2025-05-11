package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricHook struct {
	cfg               *Extra
	dbName            string
	requestsDurations *prometheus.HistogramVec
	requestsInFlight  *prometheus.GaugeVec
}

func NewMetricHook(cfg *Extra, dbName string) *metricHook {
	if cfg.Registry == nil {
		cfg.Registry = prometheus.DefaultRegisterer
	}
	registry := promauto.With(cfg.Registry)

	hook := &metricHook{
		cfg:    cfg,
		dbName: dbName,
	}

	hook.requestsDurations = registry.NewHistogramVec(prometheus.HistogramOpts{
		Name:        buildFQ("requests_durations"),
		Help:        "Requests durations",
		ConstLabels: prometheus.Labels{"db": dbName, "driver": "sqlx"},
	}, []string{"query"})

	hook.requestsInFlight = registry.NewGaugeVec(prometheus.GaugeOpts{
		Name:        buildFQ("requests_in_flight"),
		Help:        "Requests in flight",
		ConstLabels: prometheus.Labels{"db": dbName, "driver": "sqlx"},
	}, []string{"query"})

	return hook
}

func (h *metricHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	h.requestsInFlight.WithLabelValues(query).Inc()

	return context.WithValue(ctx, "begin", time.Now()), nil
}

func (h *metricHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin, ok := ctx.Value("begin").(time.Time)
	if !ok {
		return ctx, nil
	}

	go func(begin time.Time) {
		h.requestsDurations.WithLabelValues(query).Observe(time.Since(begin).Seconds())
		h.requestsInFlight.WithLabelValues(query).Dec()
	}(begin)

	return ctx, nil
}

// PgSqlxStatsCollector is a Prometheus collector for sqlx metrics.
// It implements the prometheus.Collector interface.
type PgSqlxStatsCollector struct {
	db *sqlx.DB

	acquireConns            *prometheus.Desc
	canceledAcquireCount    *prometheus.Desc
	constructingConns       *prometheus.Desc
	emptyAcquireCount       *prometheus.Desc
	idleConns               *prometheus.Desc
	maxConns                *prometheus.Desc
	totalConns              *prometheus.Desc
	newConnsCount           *prometheus.Desc
	maxLifetimeDestroyCount *prometheus.Desc
	maxIdleDestroyCount     *prometheus.Desc
}

// NewPgSqlxStatsCollector returns a new sqlxCollector.
// The dbName parameter is used to set the "db" label on the metrics.
// The db parameter is the sqlx.DB to collect metrics from.
// The db parameter must not be nil.
// The dbName parameter must not be empty.
func NewPgSqlxStatsCollector(db *sqlx.DB, dbName string) *PgSqlxStatsCollector {
	return &PgSqlxStatsCollector{
		db: db,
		acquireConns: prometheus.NewDesc(
			buildFQ("acquire_connections"),
			"Number of connections currently in the process of being acquired",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		canceledAcquireCount: prometheus.NewDesc(
			buildFQ("canceled_acquire_count"),
			"Number of times a connection acquire was canceled",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		constructingConns: prometheus.NewDesc(
			buildFQ("constructing_connections"),
			"Number of connections currently in the process of being constructed",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		emptyAcquireCount: prometheus.NewDesc(
			buildFQ("empty_acquire_count"),
			"Number of times a connection acquire was canceled",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		idleConns: prometheus.NewDesc(
			buildFQ("idle_connections"),
			"Number of idle connections in the pool",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		maxConns: prometheus.NewDesc(
			buildFQ("max_connections"),
			"Maximum number of connections allowed in the pool",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		totalConns: prometheus.NewDesc(
			buildFQ("total_connections"),
			"Total number of connections in the pool",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		newConnsCount: prometheus.NewDesc(
			buildFQ("new_connections_count"),
			"Number of new connections created",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		maxLifetimeDestroyCount: prometheus.NewDesc(
			buildFQ("max_lifetime_destroy_count"),
			"Number of connections destroyed due to MaxLifetime",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
		maxIdleDestroyCount: prometheus.NewDesc(
			buildFQ("max_idle_destroy_count"),
			"Number of connections destroyed due to MaxIdleTime",
			nil,
			prometheus.Labels{"db": dbName, "driver": "sqlx"},
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (p PgSqlxStatsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- p.acquireConns
	descs <- p.canceledAcquireCount
	descs <- p.constructingConns
	descs <- p.emptyAcquireCount
	descs <- p.idleConns
	descs <- p.maxConns
	descs <- p.totalConns
	descs <- p.newConnsCount
	descs <- p.maxLifetimeDestroyCount
	descs <- p.maxIdleDestroyCount
}

// Collect implements the prometheus.Collector interface.
func (p PgSqlxStatsCollector) Collect(metrics chan<- prometheus.Metric) {
	stats := p.db.Stats()

	metrics <- prometheus.MustNewConstMetric(p.acquireConns, prometheus.GaugeValue, float64(stats.InUse))                // Assuming InUse as acquireConns
	metrics <- prometheus.MustNewConstMetric(p.canceledAcquireCount, prometheus.CounterValue, float64(stats.WaitCount))  // Assuming WaitCount as canceledAcquireCount
	metrics <- prometheus.MustNewConstMetric(p.constructingConns, prometheus.GaugeValue, float64(stats.Idle))            // Assuming Idle as constructingConns
	metrics <- prometheus.MustNewConstMetric(p.emptyAcquireCount, prometheus.CounterValue, float64(stats.MaxIdleClosed)) // Assuming MaxIdleClosed as emptyAcquireCount
	metrics <- prometheus.MustNewConstMetric(p.idleConns, prometheus.GaugeValue, float64(stats.Idle))
	metrics <- prometheus.MustNewConstMetric(p.maxConns, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	metrics <- prometheus.MustNewConstMetric(p.totalConns, prometheus.GaugeValue, float64(stats.OpenConnections))
	metrics <- prometheus.MustNewConstMetric(p.newConnsCount, prometheus.CounterValue, float64(stats.MaxIdleTimeClosed)) // Assuming MaxIdleTimeClosed as newConnsCount
	metrics <- prometheus.MustNewConstMetric(p.maxLifetimeDestroyCount, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
	metrics <- prometheus.MustNewConstMetric(p.maxIdleDestroyCount, prometheus.CounterValue, float64(stats.WaitDuration.Seconds())) // Assuming WaitDuration as maxIdleDestroyCount
}

func buildFQ(name string) string {
	return prometheus.BuildFQName("pg", "driver", name)
}
