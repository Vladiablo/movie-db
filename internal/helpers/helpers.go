package helpers

func CopySlice[T any](slice []T) []T {
	if slice == nil {
		return nil
	}

	newSlice := make([]T, len(slice))
	copy(newSlice, slice)

	return newSlice
}
