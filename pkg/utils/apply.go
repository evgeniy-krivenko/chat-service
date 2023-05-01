package utils

func Apply[In any, Out any](in []In, f func(v In) Out) []Out {
	result := make([]Out, 0, len(in))
	for _, v := range in {
		result = append(result, f(v))
	}
	return result
}
