package routes

import (
	"context"
	"filething/models"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func GetUsers(c echo.Context) error {
	db := c.Get("db").(*bun.DB)

	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page number")
	}

	offset := page * 30
	limit := 30

	var users []models.User
	err = db.NewSelect().
		Model(&users).
		Limit(limit).
		Offset(offset).
		Scan(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(http.StatusOK, users)
}

func GetUsersCount(c echo.Context) error {
	db := c.Get("db").(*bun.DB)

	count, err := db.NewSelect().Model(&models.User{}).Count(context.Background())

	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(http.StatusOK, map[string]int{"total_users": count})
}

// Stolen from Gitea https://github.com/go-gitea/gitea
func SystemStatus(c echo.Context) error {
	updateSystemStatus()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"uptime":                     sysStatus.StartTime,
		"num_goroutine":              sysStatus.NumGoroutine,
		"cur_mem_usage":              sysStatus.MemAllocated,
		"total_mem_usage":            sysStatus.MemTotal,
		"mem_obtained":               sysStatus.MemSys,
		"ptr_lookup_times":           sysStatus.Lookups,
		"mem_allocations":            sysStatus.MemMallocs,
		"mem_frees":                  sysStatus.MemFrees,
		"cur_heap_usage":             sysStatus.HeapAlloc,
		"heap_mem_obtained":          sysStatus.HeapSys,
		"heap_mem_idle":              sysStatus.HeapIdle,
		"heap_mem_inuse":             sysStatus.HeapInuse,
		"heap_mem_release":           sysStatus.HeapReleased,
		"heap_objects":               sysStatus.HeapObjects,
		"bootstrap_stack_usage":      sysStatus.StackInuse,
		"stack_mem_obtained":         sysStatus.StackSys,
		"mspan_structures_usage":     sysStatus.MSpanInuse,
		"mspan_structures_obtained":  sysStatus.MSpanSys,
		"mcache_structures_usage":    sysStatus.MSpanInuse,
		"mcache_structures_obtained": sysStatus.MCacheSys,
		"buck_hash_sys":              sysStatus.BuckHashSys,
		"gc_sys":                     sysStatus.GCSys,
		"other_sys":                  sysStatus.OtherSys,
		"next_gc":                    sysStatus.NextGC,
		"last_gc_time":               sysStatus.LastGCTime,
		"pause_total_ns":             sysStatus.PauseTotalNs,
		"pause_ns":                   sysStatus.PauseNs,
		"num_gc":                     sysStatus.NumGC,
	})
}

var AppStartTime time.Time
var sysStatus struct {
	StartTime    string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGCTime   string // last run time
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func updateSystemStatus() {
	sysStatus.StartTime = AppStartTime.Format(time.RFC3339)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = FileSize(int64(m.Alloc))
	sysStatus.MemTotal = FileSize(int64(m.TotalAlloc))
	sysStatus.MemSys = FileSize(int64(m.Sys))
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = FileSize(int64(m.HeapAlloc))
	sysStatus.HeapSys = FileSize(int64(m.HeapSys))
	sysStatus.HeapIdle = FileSize(int64(m.HeapIdle))
	sysStatus.HeapInuse = FileSize(int64(m.HeapInuse))
	sysStatus.HeapReleased = FileSize(int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = FileSize(int64(m.StackInuse))
	sysStatus.StackSys = FileSize(int64(m.StackSys))
	sysStatus.MSpanInuse = FileSize(int64(m.MSpanInuse))
	sysStatus.MSpanSys = FileSize(int64(m.MSpanSys))
	sysStatus.MCacheInuse = FileSize(int64(m.MCacheInuse))
	sysStatus.MCacheSys = FileSize(int64(m.MCacheSys))
	sysStatus.BuckHashSys = FileSize(int64(m.BuckHashSys))
	sysStatus.GCSys = FileSize(int64(m.GCSys))
	sysStatus.OtherSys = FileSize(int64(m.OtherSys))

	sysStatus.NextGC = FileSize(int64(m.NextGC))
	sysStatus.LastGCTime = time.Unix(0, int64(m.LastGC)).Format(time.RFC3339)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
}

func FileSize(s int64) string {
	return humanize.IBytes(uint64(s))
}
