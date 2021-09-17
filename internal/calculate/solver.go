package calculate

import (
	"soultide/internal/blocks"
	doll2 "soultide/internal/doll"
	"soultide/internal/food"
	"soultide/internal/meal"
)

type solver struct {
	mealCnt    int
	index2Meal []int // 每个槽位对应的菜索引
	seeds      []int // 种子的数量
	foods      []int // 食材的数量

	coinThreshold       float64
	attractionDeviation int
	attractionTarget    int

	res Result

	p currParams
}

// SolverConfig TODO
type SolverConfig struct {
	UserFoods           map[string]int // 当前拥有的数量
	UserSeeds           map[string]int // 拥有的种子
	Name                string         // 人偶名字
	CoinThreshold       float64        // 花费的临界值
	AttractionDeviation int
	AttractionTarget    int
}

type currParams struct {
	currAttraction int
	currTimeCost   float64 // 种地的消耗时间
	currTimeWait   float64 // 商店可能的等待时间
	currTimeSeed   float64 // 种子消耗时间
	currTimeTrue   float64
	currCoinCost   float64
	currRes        []int // 当前菜，
	currNeededFood []int // 当前需要的食物
	currBlocks     []int
}

// Result TODO
type Result struct {
	Result           []int
	ResTimeCost      float64
	ResTimeCostFarm  float64 // 种地
	ResTimeCostStore float64 // 商店直接买
	ResTimeCostSeed  float64 // 买种子的消耗
	ResCoinCost      float64
	ResAttraction    int
	ResNeededFood    []int
	ResBlocks        []int
}

func (s *solver) validate() bool {
	if s.p.currAttraction > s.attractionTarget+s.attractionDeviation ||
		s.p.currAttraction < s.attractionTarget-s.attractionDeviation {
		return false
	}
	// 统计当前的需要的食物总量
	for mealIndex, mealNum := range s.p.currRes {
		if mealNum > 0 {
			currMeal := meal.GetMealById(s.index2Meal[mealIndex])
			for foodId, num := range currMeal.Foods {
				s.p.currNeededFood[foodId] += num * mealNum
				s.p.currCoinCost += food.GetFoodByFoodId(foodId).CoinCost * float64(num*mealNum)
			}
		}
	}
	// 减去已经有的
	for foodId, num := range s.p.currNeededFood {
		if num > s.foods[foodId] {
			s.p.currNeededFood[foodId] -= s.foods[foodId]
		} else {
			s.p.currNeededFood[foodId] = 0
		}
	}
	// 花费已经超过了
	if s.p.currCoinCost/float64(s.p.currAttraction) > s.coinThreshold {
		return false
	}
	// 尝试街区
	blockSolver := blocks.NewBlockSlover(s.p.currNeededFood, s.seeds)
	res := blockSolver.Solve()
	s.p.currTimeTrue = res.ResTimeCost
	s.p.currNeededFood = res.ResFoodNeeded
	s.p.currTimeWait = res.StoreTimeCost
	s.p.currTimeCost = res.FarmTimeCost
	s.p.currBlocks = res.ResResult
	s.p.currTimeSeed = res.SeedTimeCost
	s.p.currCoinCost = res.ResCoinCost
	return true
}

func (s *solver) solve(depth int) {
	if depth == s.mealCnt {
		s.resetCurrParam()
		if !s.validate() {
			return
		}
		if s.p.currTimeTrue < s.res.ResTimeCost {
			s.res.Result = NewMealSliceByCnt(s.mealCnt)
			copy(s.res.Result, s.p.currRes)
			s.res.ResTimeCost = s.p.currTimeTrue
			s.res.ResCoinCost = s.p.currCoinCost
			s.res.ResAttraction = s.p.currAttraction
			s.res.ResNeededFood = s.p.currNeededFood
			s.res.ResTimeCostFarm = s.p.currTimeCost
			s.res.ResTimeCostStore = s.p.currTimeWait
			s.res.ResTimeCostSeed = s.p.currTimeSeed
			s.res.ResBlocks = s.p.currBlocks

			//fmt.Println(s.res.ResTimeCostStore)
			//fmt.Println(s.res.ResNeededFood)
		}
		return
	}

	currMealId := s.index2Meal[depth]
	currMeal := meal.ParsedMeals[currMealId]
	maxCount := s.attractionTarget / currMeal.Attraction
	currAttraction := s.p.currAttraction
	for i := 0; i <= maxCount; i++ {
		s.p.currAttraction = currAttraction + currMeal.Attraction*i
		s.p.currRes[depth] = i
		s.solve(depth + 1)
	}

}
func (s *solver) resetCurrParam() {
	s.p.currCoinCost = 0
	s.p.currTimeCost = 0
	s.p.currNeededFood = blocks.NewFoodSlice()
	s.p.currTimeWait = 0
	s.p.currTimeTrue = 0
}

// NewMealSliceByCnt TODO
func NewMealSliceByCnt(cnt int) []int {
	return make([]int, cnt)
}

func (s *solver) reset() {
	s.res.ResTimeCost = 99999
	s.resetCurrParam()
	s.p.currRes = NewMealSliceByCnt(s.mealCnt)
}

// NewSolver TODO
func NewSolver(config SolverConfig) *solver {
	dollId := doll2.GetDollIdByname(config.Name)
	doll := doll2.GetDollById(dollId)
	s := &solver{
		mealCnt:             len(doll.Favorites),
		index2Meal:          nil,
		seeds:               nil,
		foods:               nil,
		coinThreshold:       0,
		attractionDeviation: 0,
		attractionTarget:    0,
		res:                 Result{},
		p:                   currParams{},
	}
	s.index2Meal = make([]int, s.mealCnt)
	for k, v := range doll.Favorites {
		s.index2Meal[k] = meal.GetMealIdByName(v)
	}
	s.foods = blocks.NewFoodSlice()
	for k, v := range config.UserFoods {
		s.foods[food.GetFoodIdByName(k)] = v
	}
	s.seeds = blocks.NewFoodSlice()
	for k, v := range config.UserSeeds {
		s.seeds[food.GetFoodIdByName(k)] = v
	}
	s.coinThreshold = config.CoinThreshold
	s.attractionDeviation = config.AttractionDeviation
	s.attractionTarget = config.AttractionTarget
	s.reset()
	return s
}

// Solve TODO
func Solve(config SolverConfig) Result {
	solver := NewSolver(config)
	solver.solve(0)
	return solver.res
}
