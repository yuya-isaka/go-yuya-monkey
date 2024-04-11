package object

type Environment struct {
	store map[string]Object
	outer *Environment // 拡張元の環境（外側の環境）の参照
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	// 今の環境になかったら外側の環境を調べる（再帰的に）
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	// これは変数スコープに関する考え方も定義している
	// 外側のスコープは内側のスコープを包み込む（内側になかったら外側を探す）
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
