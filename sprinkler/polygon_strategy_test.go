package sprinkler

import (
	"fmt"
	"testing"

	"github.com/zealllot/sprinklers/model"
)

// TestCase 定义测试用例结构
type TestCase struct {
	name               string
	room               *model.Room
	expectedSprinklers []model.Point
}

func TestPolygonStrategy_PlaceSprinklers(t *testing.T) {
	// 定义测试用例
	testCases := []TestCase{
		{
			name: "六边形房间",
			room: &model.Room{
				Walls: []model.Point{
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
			},
			expectedSprinklers: []model.Point{
				{X: 1500, Y: 1500}, // 左边小正方形中心
				{X: 4250, Y: 1500}, // 右边长方形第一行第一个
				{X: 6750, Y: 1500}, // 右边长方形第一行第二个
				{X: 4250, Y: 4500}, // 右边长方形第二行第一个
				{X: 6750, Y: 4500}, // 右边长方形第二行第二个
			},
		},
		{
			name: "六边形房间2",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 3000},
					{X: 3000, Y: 3000},
					{X: 3000, Y: 9000},
					{X: 9000, Y: 9000},
					{X: 9000, Y: 0},
				},
				MinSprinklerDistance:    1800,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 1700,
			},
			expectedSprinklers: []model.Point{
				{X: 1500, Y: 1500},
				{X: 4500, Y: 1500},
				{X: 7500, Y: 1500},
				{X: 4500, Y: 4500},
				{X: 7500, Y: 4500},
				{X: 4500, Y: 7500},
				{X: 7500, Y: 7500},
			},
		},
		{
			name: "六边形房间3",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 9000},
					{X: 3000, Y: 9000},
					{X: 3000, Y: 3000},
					{X: 9000, Y: 3000},
					{X: 9000, Y: 0},
				},
				MinSprinklerDistance:    1800,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 1700,
			},
			expectedSprinklers: []model.Point{
				{X: 1500, Y: 1500},
				{X: 4500, Y: 1500},
				{X: 7500, Y: 1500},
				{X: 4500, Y: 1500},
				{X: 7500, Y: 1500},
			},
		},
		{
			name: "四边形房间",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 4000},
					{X: 4000, Y: 4000},
					{X: 4000, Y: 0},
				},
				MinSprinklerDistance:    1800,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 1700,
			},
			expectedSprinklers: []model.Point{
				{X: 1000, Y: 1000},
				{X: 3000, Y: 1000},
				{X: 1000, Y: 3000},
				{X: 3000, Y: 3000},
			},
		},
		{
			name: "四边形房间2",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 3000},
					{X: 3000, Y: 3000},
					{X: 3000, Y: 0},
				},
				MinSprinklerDistance:    1800,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 1700,
			},
			expectedSprinklers: []model.Point{
				{X: 1500, Y: 1500},
			},
		},
		{
			name: "四边形房间，500喷头",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 4000},
					{X: 4000, Y: 4000},
					{X: 4000, Y: 0},
				},
				MinSprinklerDistance:    600,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 500,
			},
			expectedSprinklers: []model.Point{
				{X: 500, Y: 500},
				{X: 1500, Y: 500},
				{X: 2500, Y: 500},
				{X: 3500, Y: 500},
				{X: 500, Y: 1500},
				{X: 1500, Y: 1500},
				{X: 2500, Y: 1500},
				{X: 3500, Y: 1500},
				{X: 500, Y: 2500},
				{X: 1500, Y: 2500},
				{X: 2500, Y: 2500},
				{X: 3500, Y: 2500},
				{X: 500, Y: 3500},
				{X: 1500, Y: 3500},
				{X: 2500, Y: 3500},
				{X: 3500, Y: 3500},
			},
		},
		{
			name: "八边形房间",
			room: &model.Room{
				Walls: []model.Point{
					{X: 0, Y: 0},
					{X: 0, Y: 9000},
					{X: 3000, Y: 9000},
					{X: 3000, Y: 8000},
					{X: 2000, Y: 8000},
					{X: 2000, Y: 6000},
					{X: 9000, Y: 6000},
					{X: 9000, Y: 0},
				},
				MinSprinklerDistance:    1800,
				MinWallDistance:         100,
				SprinklerCoverageRadius: 1700,
			},
			expectedSprinklers: []model.Point{
				{X: 1500, Y: 1500},
				{X: 4500, Y: 1500},
				{X: 7500, Y: 1500},
				{X: 1500, Y: 4500},
				{X: 4500, Y: 4500},
				{X: 7500, Y: 4500},
				{X: 1500, Y: 7500},
			},
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建策略实例
			strategy := NewPolygonStrategy(tc.room, &model.Sprinkler{})

			// 获取计算结果
			result := strategy.PlaceSprinklers()

			// 验证喷头数量
			if len(result) != len(tc.expectedSprinklers) {
				t.Errorf("期望喷头数量 %d, 实际获得 %d", len(tc.expectedSprinklers), len(result))
			}

			// 验证每个喷头位置
			for i, expected := range tc.expectedSprinklers {
				found := false
				for _, actual := range result {
					// 允许有小的误差（例如0.1）
					if testDistance(expected, actual) < 0.1 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("未找到期望的喷头位置 %d: (%.2f, %.2f)", i+1, expected.X, expected.Y)
				}
			}

			// 打印实际结果，方便调试
			fmt.Println("\n实际喷头位置：")
			for i, sprinkler := range result {
				fmt.Printf("喷头 %d: (%.2f, %.2f)\n", i+1, sprinkler.X, sprinkler.Y)
			}
		})
	}
}

// testDistance 计算两点之间的距离
func testDistance(p1, p2 model.Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return dx*dx + dy*dy
}

// 您可以添加更多的测试用例，例如：
/*
func TestPolygonStrategy_CustomRoom(t *testing.T) {
	// 在这里添加您的自定义测试用例
}
*/
