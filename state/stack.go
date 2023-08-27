package state

type Object struct {
	Object  interface{}
	Order   int
	ObjType string
}

func NewObject(obj interface{}, order int, objtype string) Object {
	return Object{
		Object:  obj,
		Order:   order,
		ObjType: objtype,
	}
}

type Stack struct {
	Name                  string
	RestorationObjectName string
	Objects               map[int][]Object //int is the order in which we restore
}
