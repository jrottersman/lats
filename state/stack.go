package state

type Object struct {
	Object interface{}
	Order  int
}

type Stack struct {
	Name                  string
	RestorationObjectName string
	Objects               map[int][]Object //int is the order in which we restore
}
