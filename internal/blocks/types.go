package blocks

import (
	"math"
	"soultide/internal/food"
)

var (
	// SimplifiedBlocks TODO
	SimplifiedBlocks = make(map[string]block)
)

func init() {
	for _, v := range primitiveBlocks {
		SimplifiedBlocks[v.name] = v
	}
}

type block struct {
	name         string
	acquireFoods []string
}

// BlockSolver TODO
type BlockSolver struct {
	currParam currParam
	blocks    []block

	BlockResult

	OriginFoodNeeded map[string]int
	OriginSeedOwned  map[string]int
}
type currParam struct {
	fre        int
	currNeeded map[string]int
	currResult map[string]int
}

// BlockResult TODO
type BlockResult struct {
	StoreTimeCost float64
	FarmTimeCost  float64
	SeedTimeCost  float64
	ResTimeCost   float64
	ResFoodNeeded map[string]int
	ResResult     map[string]int
}

// NewBlockSlover TODO
func NewBlockSlover(foodNeeded map[string]int, seedsOwned map[string]int) *BlockSolver {
	bs := &BlockSolver{
		currParam: currParam{},
		blocks:    primitiveBlocks,
		BlockResult: BlockResult{
			ResTimeCost:   9999,
			ResFoodNeeded: nil,
			ResResult:     nil,
		},
		OriginFoodNeeded: foodNeeded,
		OriginSeedOwned:  seedsOwned,
	}
	bs.resetCurrParam()
	return bs
}

func (s *BlockSolver) resetCurrParam() {
	s.currParam.reset(s.OriginFoodNeeded)
}

// Solve TODO
func (s *BlockSolver) Solve() BlockResult {
	s.solve(0)
	return s.BlockResult
}

func (s *BlockSolver) solve(depth int) {
	if depth == len(s.blocks) {
		s.currParam.resetNeededMap(s.OriginFoodNeeded)
		// 巡查完成
		// 先计算巡查可以减去的数量
		//	fmt.Println("--------------------------")
		for k, v := range s.currParam.currResult {
			for _, fn := range SimplifiedBlocks[k].acquireFoods {
				if s.currParam.currNeeded[fn] != 0 {
					if s.currParam.currNeeded[fn] >= 3*v {
						s.currParam.currNeeded[fn] -= 3 * v
					} else {
						s.currParam.currNeeded[fn] = 0
					}
				}
			}
		}

		// 分别计算购物所需的时间和种地所需时间
		// 购物所需时间取得最大值，种地的话假如每天只种16块地
		// 加入种子消耗
		var timeCost float64
		var waitCost float64
		var seedCost float64
		for k, v := range s.currParam.currNeeded {
			if food.SimpifiedFood[k].TimeCost == 0 {
				if waitCost < float64(v)/20*24 {
					waitCost = float64(v) / 20 * 24
				}
			} else {
				// 种地消耗计算
				if v%5 != 0 {
					timeCost += 1
				}
				timeCost += float64(v / 5)
				// 种子消耗计算
				seedOwned := 0
				if v, ok := s.OriginSeedOwned[k]; ok {
					seedOwned = v
				}
				if v > seedOwned*5 {
					v = v - seedOwned*5
				} else {
					v = 0
				}
				if seedCost < float64(v)/25*24 {
					seedCost = float64(v) / 25 * 24
				}
			}
		}
		timeCost = timeCost / 16 * 24
		trueCost := math.Max(timeCost, waitCost)
		trueCost = math.Max(trueCost, seedCost)
		if s.ResTimeCost > trueCost {
			s.StoreTimeCost = waitCost
			s.FarmTimeCost = timeCost
			s.SeedTimeCost = seedCost
			s.ResTimeCost = trueCost
			s.ResResult = make(map[string]int)
			s.ResFoodNeeded = make(map[string]int)
			for k, v := range s.currParam.currResult {
				s.ResResult[k] = v
			}
			for k, v := range s.currParam.currNeeded {
				s.ResFoodNeeded[k] = v
			}
		}
		return
	}
	currBlock := s.blocks[depth]
	currFre := s.currParam.fre
	// 减去获得
	for i := 0; i <= 2; i++ {
		if currFre+i > 5 {
			break
		}
		s.currParam.fre = currFre + i
		s.currParam.currResult[currBlock.name] = i
		//	fmt.Println(currBlock.name, i)
		s.solve(depth + 1)
	}
}

func (curr *currParam) resetNeededMap(needed map[string]int) {
	for k, v := range needed {
		curr.currNeeded[k] = v
	}
}

func (curr *currParam) reset(needed map[string]int) {
	curr.fre = 0
	curr.currNeeded = make(map[string]int)
	curr.resetNeededMap(needed)
	curr.currResult = make(map[string]int)
}

// TryToPatrol TODO
func TryToPatrol(needed map[string]int) {

}
