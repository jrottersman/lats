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

func NewStack(name string, restorationObjectName string, objects []Object) Stack {
	objs := make(map[int][]Object)
	for _, v := range objects {
		order := v.Order
		_, ok := objs[order]
		if !ok {
			objs[order] = []Object{v}
		}
	}

	return Stack{
		Name:                  name,
		RestorationObjectName: restorationObjectName,
		Objects:               objs,
	}
}
