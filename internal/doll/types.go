package doll

var (
	// SimplifiedDolls TODO
	SimplifiedDolls = make(map[string]int)
)

// GetDollIdByname TODO
func GetDollIdByname(name string) int {
	return SimplifiedDolls[name]
}

// GetDollById TODO
func GetDollById(id int) Doll {
	return dolls[id]
}

// Doll TODO
type Doll struct {
	name      string
	Favorites []string
}

// InitDolls TODO
func InitDolls() {
	for k, v := range dolls {
		SimplifiedDolls[v.name] = k
	}
}
