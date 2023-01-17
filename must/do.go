package must

func Do(err error) {
	if err != nil {
		panic(err)
	}
}

func Do2[T any](v T, err error) T {
	Do(err)
	return v
}
