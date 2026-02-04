package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/moganji/jot/pkg/models"
	"github.com/moganji/jot/pkg/storage"
)

type BenchResult struct {
	Notes      int
	Categories int
	Words      int
	AddAvg     float64
	SearchAvg  float64
	ListAvg    float64
	GetAvg     float64
	UpdateAvg  float64
	DeleteAvg  float64
}

var words = []string{"api", "key", "stripe", "payment", "webhook", "signature", "verify", "endpoint", "request", "response", "error", "success", "data", "user", "authentication", "authorization", "token", "session", "database", "query", "index", "search", "cache", "redis", "postgres", "mongodb", "mysql", "docker", "kubernetes", "microservice", "lambda", "function", "server", "client", "network", "protocol", "http", "https", "rest", "graphql", "json", "xml", "configuration", "environment", "production", "staging", "development", "deployment"}

func generateContent(wordCount int) string {
	content := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		content[i] = words[rand.Intn(len(words))]
	}
	return strings.Join(content, " ")
}

func benchmark(noteCount, categoryCount, wordCount int) BenchResult {
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("jot-bench-%d", time.Now().UnixNano()))
	defer os.RemoveAll(tmpDir)

	store, _ := storage.NewFileStore(tmpDir)
	defer store.Close()

	categories := make([]string, categoryCount)
	for i := 0; i < categoryCount; i++ {
		categories[i] = fmt.Sprintf("cat-%d", i)
	}

	noteIDs := make([]string, noteCount)
	content := generateContent(wordCount)

	start := time.Now()
	for i := 0; i < noteCount; i++ {
		cat := categories[i%categoryCount]
		note := models.NewNote(cat, fmt.Sprintf("note-%d", i), content, nil)
		noteIDs[i] = note.ID
		store.Add(note)
	}
	addTime := time.Since(start)
	addAvg := float64(addTime.Microseconds()) / float64(noteCount) / 1000.0

	searchQueries := []string{"api", "stripe", "database", "authentication", "error"}
	start = time.Now()
	for _, q := range searchQueries {
		store.Search(q, "", nil)
	}
	searchTime := time.Since(start)
	searchAvg := float64(searchTime.Microseconds()) / float64(len(searchQueries)) / 1000.0

	testCategories := min(10, categoryCount)
	start = time.Now()
	for i := 0; i < testCategories; i++ {
		store.List(categories[i])
	}
	listTime := time.Since(start)
	listAvg := float64(listTime.Microseconds()) / float64(testCategories) / 1000.0

	testGets := min(100, noteCount)
	start = time.Now()
	for i := 0; i < testGets; i++ {
		store.Get(noteIDs[rand.Intn(noteCount)])
	}
	getTime := time.Since(start)
	getAvg := float64(getTime.Microseconds()) / float64(testGets) / 1000.0

	testUpdates := min(50, noteCount)
	start = time.Now()
	for i := 0; i < testUpdates; i++ {
		note, _ := store.Get(noteIDs[rand.Intn(noteCount)])
		if note != nil {
			note.Content = note.Content + " updated"
			store.Update(note)
		}
	}
	updateTime := time.Since(start)
	updateAvg := float64(updateTime.Microseconds()) / float64(testUpdates) / 1000.0

	testDeletes := min(50, noteCount)
	start = time.Now()
	for i := 0; i < testDeletes; i++ {
		store.Delete(noteIDs[i])
	}
	deleteTime := time.Since(start)
	deleteAvg := float64(deleteTime.Microseconds()) / float64(testDeletes) / 1000.0

	return BenchResult{
		Notes:      noteCount,
		Categories: categoryCount,
		Words:      wordCount,
		AddAvg:     addAvg,
		SearchAvg:  searchAvg,
		ListAvg:    listAvg,
		GetAvg:     getAvg,
		UpdateAvg:  updateAvg,
		DeleteAvg:  deleteAvg,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())

	configs := []struct {
		notes      int
		categories int
		words      int
	}{
		{100, 10, 10},
		{100, 10, 100},
		{100, 10, 1000},
		{1000, 10, 100},
		{1000, 100, 100},
		{1000, 100, 1000},
		{10000, 100, 100},
		{10000, 1000, 100},
		{10000, 1000, 1000},
		{100000, 1000, 100},
		{100000, 1000, 1000},
	}

	fmt.Println("# Benchmark Results")
	fmt.Println()
	fmt.Println("| Notes | Categories | Words/Note | Add (ms) | Search (ms) | List (ms) | Get (ms) | Update (ms) | Delete (ms) |")
	fmt.Println("|-------|------------|------------|----------|-------------|-----------|----------|-------------|-------------|")

	for i, cfg := range configs {
		fmt.Fprintf(os.Stderr, "[%d/%d] Testing: %d notes, %d categories, %d words\n", i+1, len(configs), cfg.notes, cfg.categories, cfg.words)
		result := benchmark(cfg.notes, cfg.categories, cfg.words)

		fmt.Printf("| %d | %d | %d | %.3f | %.3f | %.3f | %.3f | %.3f | %.3f |\n",
			result.Notes, result.Categories, result.Words,
			result.AddAvg, result.SearchAvg, result.ListAvg,
			result.GetAvg, result.UpdateAvg, result.DeleteAvg)

		fmt.Fprintf(os.Stderr, "[%d/%d] Completed: Add=%.3fms Search=%.3fms List=%.3fms Get=%.3fms Update=%.3fms Delete=%.3fms\n",
			i+1, len(configs), result.AddAvg, result.SearchAvg, result.ListAvg, result.GetAvg, result.UpdateAvg, result.DeleteAvg)
	}
}
