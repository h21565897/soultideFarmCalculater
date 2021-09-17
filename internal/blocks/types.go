package blocks

import (
	"math"
	"soultide/internal/food"
)

var (
	// ParsedBlocks TODO
	// SimplifiedBlocks = make(map[string]block)
	ParsedBlocks = make([]nblock, len(primitiveBlocks))

	// BlockCnt TODO
	BlockCnt int
)

// GetBlockByBlockId TODO
func GetBlockByBlockId(id int) nblock {
	return ParsedBlocks[id]
}

// InitBlock TODO
func InitBlock() {
	for k, v := range primitiveBlocks {
		var b nblock
		b.name = v.name
		b.acquireFoods = make([]int, len(v.acquireFoods))
		for n, m := range v.acquireFoods {
			b.acquireFoods[n] = food.GetFoodIdByName(m)
		}
		ParsedBlocks[k] = b
	}
	BlockCnt = len(ParsedBlocks)
}

type block struct {
	name         string
	acquireFoods []string
}
type nblock struct {
	name         string
	acquireFoods []int // 可以获得的食物的索引
}

// BlockSolver TODO
type BlockSolver struct {
	currParam currParam

	BlockResult

	OriginFoodNeeded []int
	OriginSeedOwned  []int
}
type currParam struct {
	fre          int
	currNeeded   []int
	currResult   []int
	currCoinCost float64
}

// BlockResult TODO
type BlockResult struct {
	StoreTimeCost float64
	FarmTimeCost  float64
	SeedTimeCost  float64
	ResTimeCost   float64
	ResFoodNeeded []int
	ResCoinCost   float64
	ResResult     []int
}

// NewBlockSlover TODO
func NewBlockSlover(foodNeeded []int, seedsOwned []int) *BlockSolver {
	bs := &BlockSolver{
		currParam: currParam{},
		BlockResult: BlockResult{
			ResTimeCost:   9999,
			ResFoodNeeded: nil,
			ResResult:     nil,
		},
		OriginFoodNeeded: foodNeeded,
		OriginSeedOwned:  seedsOwned,
	}
	bs.currParam.reset()
	return bs
}

// Solve TODO
func (s *BlockSolver) Solve() BlockResult {
	s.solve(0)
	return s.BlockResult
}

func (s *BlockSolver) solve(depth int) {
	if depth == BlockCnt {
		s.currParam.currNeeded = NewFoodSlice()
		s.currParam.currCoinCost = 0
		copy(s.currParam.currNeeded, s.OriginFoodNeeded)
		// 巡查完成
		// 先计算巡查可以减去的数量
		//	fmt.Println("--------------------------")
		for blockid, v := range s.currParam.currResult {
			currBlock := GetBlockByBlockId(blockid)
			for _, foodId := range currBlock.acquireFoods {
				if s.currParam.currNeeded[foodId] > v*3 {
					s.currParam.currNeeded[foodId] -= v * 3
				} else {
					s.currParam.currNeeded[foodId] = 0
				}
			}
		}

		// 分别计算购物所需的时间和种地所需时间
		// 购物所需时间取得最大值，种地的话假如每天只种16块地
		// 加入种子消耗
		var timeCost float64
		var waitCost float64
		var seedCost float64
		for foodId, v := range s.currParam.currNeeded {
			currFood := food.GetFoodByFoodId(foodId)
			if currFood.TimeCost == 0 {
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
				seedOwned := s.OriginSeedOwned[foodId]
				if v > seedOwned*5 {
					v = v - seedOwned*5
				} else {
					v = 0
				}
				if seedCost < float64(v)/25*24 {
					seedCost = float64(v) / 25 * 24
				}
			}
			s.currParam.currCoinCost += food.GetFoodByFoodId(foodId).CoinCost * float64(v)
		}
		timeCost = timeCost / 16 * 24
		trueCost := math.Max(timeCost, waitCost)
		trueCost = math.Max(trueCost, seedCost)
		if s.ResTimeCost > trueCost {
			s.StoreTimeCost = waitCost
			s.FarmTimeCost = timeCost
			s.SeedTimeCost = seedCost
			s.ResTimeCost = trueCost
			s.ResResult = NewBlockSlice()
			s.ResFoodNeeded = NewFoodSlice()
			s.ResCoinCost = s.currParam.currCoinCost
			copy(s.ResResult, s.currParam.currResult)
			copy(s.ResFoodNeeded, s.currParam.currNeeded)
		}
		return
	}
	currFre := s.currParam.fre
	// 减去获得
	for i := 0; i <= 2; i++ {
		if currFre+i > 5 {
			break
		}
		s.currParam.fre = currFre + i
		s.currParam.currResult[depth] = i
		//	fmt.Println(currBlock.name, i)
		s.solve(depth + 1)
	}
}

// NewFoodSlice TODO
func NewFoodSlice() []int {
	return make([]int, food.FoodCnt)
}

// NewBlockSlice TODO
func NewBlockSlice() []int {
	return make([]int, BlockCnt)
}

func (curr *currParam) reset() {
	curr.fre = 0
	curr.currResult = NewBlockSlice()
}

// TryToPatrol TODO
func TryToPatrol(needed map[string]int) {

}
