package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// プレイヤーのセーブ状態の構造体
type SaveState struct {
	CurrentNode string          `json:"current_node"`
	Flags       map[string]bool `json:"flags"`
}

var (
	scenarioCache []byte
	saveStore     = make(map[string]SaveState) // 本来はRedisやDBに入れる領域
	mu            sync.RWMutex
)

func main() {
	// 1. PHPが生成したシナリオJSONをメモリにロード（ディスクI/Oを殺す）
	data, err := ioutil.ReadFile("scenario.json")
	if err != nil {
		log.Fatalf("シナリオファイルの読み込み失敗: %v", err)
	}
	scenarioCache = data

	// 2. ルーティング設定（CORS対策込み）
	http.HandleFunc("/api/scenario", handleScenario)
	http.HandleFunc("/api/save", handleSave)

	log.Println("[Go] API Server running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// シナリオデータを一瞬で返すAPI
func handleScenario(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" { return }
	w.Header().Set("Content-Type", "application/json")
	w.Write(scenarioCache)
}

// セーブデータをオンメモリで処理するAPI
func handleSave(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" { return }

	if r.Method == http.MethodPost {
		var state SaveState
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		saveStore["default_user"] = state // 簡易的に固定ユーザーで保存
		mu.Unlock()
		w.Write([]byte(`{"status":"success"}`))
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		mu.RLock()
		state, exists := saveStore["default_user"]
		mu.RUnlock()
		if !exists {
			state = SaveState{CurrentNode: "start", Flags: make(map[string]bool)}
		}
		json.NewEncoder(w).Encode(state)
	}
}

func setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}