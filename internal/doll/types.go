package doll

var (
	// SimplifiedDolls TODO
	SimplifiedDolls = make(map[string]Doll)
)

// Doll TODO
type Doll struct {
	name      string
	Favorites []string
}

// InitDolls TODO
func InitDolls() {
	for _, v := range dolls {
		SimplifiedDolls[v.name] = v
	}
}
