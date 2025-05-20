package sprinkler

import (
	"math"

	"github.com/zealllot/sprinklers/model"
)

type PolygonStrategy struct {
	room      *model.Room
	sprinkler *model.Sprinkler
}

func NewPolygonStrategy(room *model.Room, sprinkler *model.Sprinkler) *PolygonStrategy {
	return &PolygonStrategy{
		room:      room,
		sprinkler: sprinkler,
	}
}

func (s *PolygonStrategy) PlaceSprinklers() []model.Point {
	// 计算房间的边界框
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64
	for _, wall := range s.room.Walls {
		minX = math.Min(minX, wall.X)
		minY = math.Min(minY, wall.Y)
		maxX = math.Max(maxX, wall.X)
		maxY = math.Max(maxY, wall.Y)
	}

	coverage := s.room.SprinklerCoverageRadius
	step := coverage * 2

	// 尝试不同的起点偏移
	bestSprinklers := []model.Point{}
	maxCoverage := 0.0
	minSprinklers := math.MaxInt32

	// 在[0, 2R)范围内每隔R/4尝试一次
	offsetStep := step / 8
	for offsetX := 0.0; offsetX < step; offsetX += offsetStep {
		for offsetY := 0.0; offsetY < step; offsetY += offsetStep {
			// 计算这个偏移下的网格起点和终点
			startX := math.Floor((minX-offsetX)/step)*step + offsetX + coverage
			startY := math.Floor((minY-offsetY)/step)*step + offsetY + coverage
			endX := math.Ceil((maxX-offsetX)/step)*step + offsetX
			endY := math.Ceil((maxY-offsetY)/step)*step + offsetY

			// 生成当前偏移下的喷头方案
			var currentSprinklers []model.Point
			for x := startX; x <= endX; x += step {
				for y := startY; y <= endY; y += step {
					if s.isPointInPolygon(x, y) && s.checkWallDistance(x, y) {
						currentSprinklers = append(currentSprinklers, model.Point{X: x, Y: y})
					}
				}
			}

			// 评估当前方案的覆盖效果
			currentCoverage := s.evaluateCoverage(currentSprinklers)

			// 选择覆盖率最高且喷头数量最少的方案
			if currentCoverage > maxCoverage ||
				(math.Abs(currentCoverage-maxCoverage) < 1e-6 && len(currentSprinklers) < minSprinklers) {
				maxCoverage = currentCoverage
				minSprinklers = len(currentSprinklers)
				bestSprinklers = currentSprinklers
			}
		}
	}

	return bestSprinklers
}

// 评估喷头方案的覆盖效果
func (s *PolygonStrategy) evaluateCoverage(sprinklers []model.Point) float64 {
	// 计算房间边界框
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64
	for _, wall := range s.room.Walls {
		minX = math.Min(minX, wall.X)
		minY = math.Min(minY, wall.Y)
		maxX = math.Max(maxX, wall.X)
		maxY = math.Max(maxY, wall.Y)
	}

	// 采样检查覆盖情况
	const grid = 5.0 // 更细致的采样
	totalPoints := 0
	coveredPoints := 0
	coverage := s.room.SprinklerCoverageRadius

	for x := minX; x <= maxX; x += grid {
		for y := minY; y <= maxY; y += grid {
			if !s.isPointInPolygon(x, y) {
				continue
			}
			totalPoints++
			for _, sp := range sprinklers {
				if math.Abs(x-sp.X) <= coverage && math.Abs(y-sp.Y) <= coverage {
					coveredPoints++
					break
				}
			}
		}
	}

	if totalPoints == 0 {
		return 0
	}
	return float64(coveredPoints) / float64(totalPoints)
}

// 判断所有房间区域是否被完全覆盖
func (s *PolygonStrategy) isFullCovered(sprinklers []model.Point, minX, minY, maxX, maxY, coverage float64) bool {
	const grid = 10.0 // 采样精度更高
	for x := minX; x <= maxX; x += grid {
		for y := minY; y <= maxY; y += grid {
			if !s.isPointInPolygon(x, y) {
				continue
			}
			covered := false
			for _, sp := range sprinklers {
				if math.Abs(x-sp.X) <= coverage && math.Abs(y-sp.Y) <= coverage {
					covered = true
					break
				}
			}
			if !covered {
				return false
			}
		}
	}
	return true
}

func (s *PolygonStrategy) calculateBoundingBox() (minX, minY, maxX, maxY float64) {
	minX = s.room.Walls[0].X
	minY = s.room.Walls[0].Y
	maxX = s.room.Walls[0].X
	maxY = s.room.Walls[0].Y

	for _, point := range s.room.Walls {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return
}

func (s *PolygonStrategy) isValidPosition(x, y float64) bool {
	// 检查点是否在房间内
	if !s.isPointInPolygon(x, y) {
		return false
	}

	// 检查与墙体的距离
	if !s.checkWallDistance(x, y) {
		return false
	}

	// 检查正方形的四个角是否都在房间内
	halfSize := s.room.SprinklerCoverageRadius
	if halfSize == 0 {
		return true
	}

	// 检查覆盖范围是否在房间内
	// 使用更宽松的检查：只要覆盖范围的中心点在房间内，且与墙体的距离足够，就认为是有效的
	return true
}

func (s *PolygonStrategy) isPointInPolygon(x, y float64) bool {
	// 先判断是否在多边形的某条边上
	const epsilon = 1e-6
	for i := 0; i < len(s.room.Walls); i++ {
		j := (i + 1) % len(s.room.Walls)
		x1, y1 := s.room.Walls[i].X, s.room.Walls[i].Y
		x2, y2 := s.room.Walls[j].X, s.room.Walls[j].Y
		// 判断点是否在线段上
		dx := x2 - x1
		dy := y2 - y1
		if math.Abs(dx) < epsilon && math.Abs(dy) < epsilon {
			continue // 跳过零长度边
		}
		t := ((x-x1)*dx + (y-y1)*dy) / (dx*dx + dy*dy)
		if t >= -epsilon && t <= 1+epsilon {
			px := x1 + t*dx
			py := y1 + t*dy
			if math.Hypot(px-x, py-y) < epsilon {
				return true
			}
		}
	}
	// 射线法判断是否在多边形内部
	inside := false
	j := len(s.room.Walls) - 1
	for i := 0; i < len(s.room.Walls); i++ {
		xi, yi := s.room.Walls[i].X, s.room.Walls[i].Y
		xj, yj := s.room.Walls[j].X, s.room.Walls[j].Y
		if ((yi > y) != (yj > y)) &&
			x < (xj-xi)*(y-yi)/(yj-yi+epsilon)+xi {
			inside = !inside
		}
		j = i
	}
	return inside
}

func (s *PolygonStrategy) isPointOnLine(x, y, x1, y1, x2, y2 float64) bool {
	// 检查点是否在线段上
	const epsilon = 1e-6
	if math.Abs((y2-y1)*(x-x1)-(y-y1)*(x2-x1)) > epsilon {
		return false
	}
	if x < math.Min(x1, x2)-epsilon || x > math.Max(x1, x2)+epsilon {
		return false
	}
	if y < math.Min(y1, y2)-epsilon || y > math.Max(y1, y2)+epsilon {
		return false
	}
	return true
}

func (s *PolygonStrategy) checkWallDistance(x, y float64) bool {
	// 如果最小墙体距离为0，则不检查
	if s.room.MinWallDistance <= 0 {
		return true
	}

	for i := 0; i < len(s.room.Walls); i++ {
		j := (i + 1) % len(s.room.Walls)
		wallStart := s.room.Walls[i]
		wallEnd := s.room.Walls[j]

		// 计算点到线段的距离
		distance := s.pointToLineDistance(x, y, wallStart.X, wallStart.Y, wallEnd.X, wallEnd.Y)
		if distance < s.room.MinWallDistance {
			return false
		}
	}

	return true
}

func (s *PolygonStrategy) pointToLineDistance(x, y, x1, y1, x2, y2 float64) float64 {
	// 计算点到线段的距离
	A := x - x1
	B := y - y1
	C := x2 - x1
	D := y2 - y1

	dot := A*C + B*D
	lenSq := C*C + D*D
	var param float64

	if lenSq != 0 {
		param = dot / lenSq
	}

	var xx, yy float64

	if param < 0 {
		xx = x1
		yy = y1
	} else if param > 1 {
		xx = x2
		yy = y2
	} else {
		xx = x1 + param*C
		yy = y1 + param*D
	}

	dx := x - xx
	dy := y - yy

	return math.Sqrt(dx*dx + dy*dy)
}
