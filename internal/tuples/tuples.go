package tuples

type Double[TFirst any, TSecond any] struct {
	First  *TFirst
	Second *TSecond
}

type Triple[TFirst any, TSecond any, TThird any] struct {
	First  *TFirst
	Second *TSecond
	Third  *TThird
}

func (t *Double[TFirst, TSecond]) Deconstruct() (*TFirst, *TSecond) {
	return t.First, t.Second
}

func (t *Triple[TFirst, TSecond, TThird]) Deconstruct() (*TFirst, *TSecond, *TThird) {
	return t.First, t.Second, t.Third
}
