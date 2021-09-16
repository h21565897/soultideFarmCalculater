package main

import (
	"fmt"
	"soultide/internal/calculate"
	"soultide/internal/doll"
	"soultide/internal/food"
	"soultide/internal/meal"
)

func init() {
	food.InitFood()
	meal.InitMeals()
	doll.InitDolls()
}

func main() {
	uf := map[string]int{
		"面粉":  26,
		"苹果":  12,
		"西瓜":  15,
		"青菇":  14,
		"米":   6,
		"草莓":  2,
		"卷心菜": 15,
		"松茸":  8,
		"里脊":  3,
		"苍鱼":  21,
		"蜂蜜":  13,
		"牛奶":  9,
		"牛排":  14,
		"三文鱼": 9,
		"辣椒":  25,
		"冰块":  12,
	}
	dollName := "尼柯莱特"
	s := calculate.Solve(dollName, 5000, 50, 100, uf)
	fmt.Println("当前人偶：", dollName)
	fmt.Println("当前需要食物:")
	fmt.Println("共获得好感:", s.ResAttraction)
	for k, v := range s.Result {
		fmt.Println(k, ":", v, "份")
	}
	fmt.Println("推荐巡查的街区:")
	for k, v := range s.ResBlocks {
		fmt.Println(k, v, "次")
	}
	fmt.Println("当前每点好感度消耗金币数:", s.ResCoinCost/float64(s.ResAttraction), "共消耗金币数:", s.ResCoinCost)
	fmt.Println("当前每点好感度消耗小时数:", s.ResTimeCost/float64(s.ResAttraction), "共消耗小时数:", s.ResTimeCost)
	fmt.Println("其中种地需要的小时数：", s.ResTimeCostFarm, "块数:", s.ResTimeCostFarm/24*16)
	fmt.Println("其中商店购买需要等待的小时数：", s.ResTimeCostStore)
	fmt.Println("当前需求食材:")
	for k, v := range s.ResNeededFood {
		if v == 0 {
			fmt.Println(k, ":", v, "份", "可通过已经有的和巡查满足要求")
			continue
		}
		if food.SimpifiedFood[k].TimeCost == 0 {
			fmt.Println(k, ":", v, "份", "该食材在商店购买")
		} else {
			fmt.Println(k, ":", v, "份", "该食材种地获得，需要种：", v/5+1, "块地，请合理规划")
		}
	}

}
