package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"soultide/internal/blocks"
	"soultide/internal/calculate"
	"soultide/internal/doll"
	"soultide/internal/food"
	"soultide/internal/meal"
	"time"
)

func init() {
	food.InitFood()
	meal.InitMeals()
	doll.InitDolls()
	blocks.InitBlock()
}

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:8899", nil)
	}()
	uf := map[string]int{
		"面粉":  0,
		"苹果":  0,
		"西瓜":  0,
		"青菇":  0,
		"米":   0,
		"草莓":  0,
		"卷心菜": 0,
		"松茸":  0,
		"里脊":  0,
		"苍鱼":  0,
		"蜂蜜":  0,
		"牛奶":  0,
		"牛排":  0,
		"三文鱼": 0,
		"辣椒":  0,
		"冰块":  0,
	}
	us := map[string]int{
		"面粉":  0,
		"苹果":  0,
		"西瓜":  0,
		"青菇":  0,
		"米":   0,
		"草莓":  0,
		"卷心菜": 0,
		"松茸":  0,
	}
	tm := time.Now()
	dollName := "柯露雪儿"
	config := calculate.SolverConfig{
		UserFoods:           uf,
		UserSeeds:           us,
		Name:                dollName,
		CoinThreshold:       5000,
		AttractionDeviation: 200,
		AttractionTarget:    2000,
	}
	s := calculate.Solve(config)
	fmt.Println("当前人偶：", dollName)
	fmt.Println("当前需要食物:")
	fmt.Println("共获得好感:", s.ResAttraction)
	dollId := doll.GetDollIdByname(dollName)
	for k, v := range s.Result {
		fmt.Println(doll.GetDollById(dollId).Favorites[k], ":", v, "份")
	}
	fmt.Println("推荐巡查的街区:")
	for k, v := range s.ResBlocks {
		fmt.Println(blocks.GetBlockByBlockId(k), v, "次")
	}
	fmt.Println("当前每点好感度消耗金币数:", s.ResCoinCost/float64(s.ResAttraction), "共消耗金币数:", s.ResCoinCost)
	fmt.Println("当前每点好感度消耗小时数:", s.ResTimeCost/float64(s.ResAttraction), "共消耗小时数:", s.ResTimeCost)
	fmt.Println("其中种地需要的小时数：", s.ResTimeCostFarm, "块数:", s.ResTimeCostFarm/24*16)
	fmt.Println("其中商店购买需要等待的小时数：", s.ResTimeCostStore)
	fmt.Println("其中商店购买种子需要等待的小时数：", s.ResTimeCostSeed)
	fmt.Println("当前需求食材:")
	for k, v := range s.ResNeededFood {
		if v == 0 {
			fmt.Println(food.GetFoodByFoodId(k).Name, ":", v, "份", "可通过已经有的和巡查满足要求")
			continue
		}
		if food.GetFoodByFoodId(k).TimeCost == 0 {
			fmt.Println(food.GetFoodByFoodId(k).Name, ":", v, "份", "该食材在商店购买")
		} else {
			fmt.Println(food.GetFoodByFoodId(k).Name, ":", v, "份", "该食材种地获得，需要种：", v/5+1, "块地，请合理规划")
		}
	}
	fmt.Println("当前需要种子:")
	for k, v := range s.ResNeededFood {
		if food.GetFoodByFoodId(k).TimeCost != 0 {
			if v > us[food.GetFoodByFoodId(k).Name]*5 {
				fmt.Println("当前需要", food.GetFoodByFoodId(k).Name, "种子", (v-us[food.GetFoodByFoodId(k).Name]*5)/5+1, "份")
			}
		}
	}
	fmt.Println("耗时：", time.Since(tm))
}
