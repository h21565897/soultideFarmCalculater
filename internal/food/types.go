package food

const (
	// StorePrice TODO
	StorePrice = 2000
	// SeedPrice TODO
	SeedPrice = 1200
)

var (
	// SimpifiedFood TODO
	SimpifiedFood = make(map[string]Food)
)

type food struct {
	name    string
	inSotre bool // 是否是商店里面可以购买，如果是就是可以买的，不是就是不可以买的
}

// Food TODO
// 解析之后的食材（计算时间消耗和金币消耗）
type Food struct {
	Name     string
	TimeCost float64
	CoinCost float64
}

// InitFood TODO
func InitFood() {
	parseFood(primitiveFood)
}

func parseFood(foods []food) {
	for _, v := range foods {
		var nFood Food
		if v.inSotre {
			nFood = Food{
				Name:     v.name,
				TimeCost: 0,
				CoinCost: StorePrice,
			}
		} else {
			nFood = Food{
				Name:     v.name,
				TimeCost: 4.0 / 5.0,
				CoinCost: SeedPrice / 5,
			}
		}
		SimpifiedFood[nFood.Name] = nFood
	}
}
