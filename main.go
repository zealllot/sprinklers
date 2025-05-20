package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/zealllot/sprinklers/model"
	sp "github.com/zealllot/sprinklers/sprinkler"
)

type CalculationRequest struct {
	Points          []model.Point `json:"points"`
	CoverageRadius  float64       `json:"coverageRadius"`
	MinWallDistance float64       `json:"minWallDistance"`
}

type CalculationResponse struct {
	Room           *model.Room   `json:"room"`
	Sprinklers     []model.Point `json:"sprinklers"`
	CoverageRadius float64       `json:"coverageRadius"`
}

// 创建一个示例房间
var room = &model.Room{
	Walls: []model.Point{
		// A（0,0）、B（0,3000）、C（3000,3000）、D（3000,9000）、E（9000,9000）、F（9000,0）
		{X: 0, Y: 0},
		{X: 0, Y: 3000},
		{X: 3000, Y: 3000},
		{X: 3000, Y: 6000},
		{X: 8000, Y: 6000},
		{X: 8000, Y: 0},
	},
	MinSprinklerDistance:    1800,
	MinWallDistance:         100,
	SprinklerCoverageRadius: 1700,
}

func main() {
	// 设置静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 主页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(filepath.Join("sprinkler", "template.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 计算接口
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CalculationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("解析请求失败: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("收到计算请求: %+v", req)

		// 验证输入
		if len(req.Points) < 3 {
			msg := "至少需要3个端点才能形成一个房间"
			log.Printf("验证失败: %s", msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// 创建房间
		room := &model.Room{
			Walls:                   req.Points,
			MinSprinklerDistance:    req.CoverageRadius * 1.1, // 喷头之间最小距离为覆盖半径的1.1倍
			MinWallDistance:         req.MinWallDistance,
			SprinklerCoverageRadius: req.CoverageRadius,
		}

		// 创建喷头
		sprinkler := &model.Sprinkler{
			Coverage: req.CoverageRadius,
		}

		// 计算喷头位置
		strategy := sp.NewPolygonStrategy(room, sprinkler)
		sprinklers := strategy.PlaceSprinklers()

		log.Printf("计算完成，找到 %d 个喷头位置", len(sprinklers))

		// 返回结果
		response := CalculationResponse{
			Room: &model.Room{
				Walls:                   room.Walls,
				SprinklerCoverageRadius: room.SprinklerCoverageRadius,
				MinSprinklerDistance:    room.MinSprinklerDistance,
				MinWallDistance:         room.MinWallDistance,
			},
			Sprinklers:     sprinklers,
			CoverageRadius: req.CoverageRadius,
		}

		log.Printf("返回响应: %+v", response)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("编码响应失败: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("服务器启动在 http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
