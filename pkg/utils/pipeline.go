package utils

func addSorter(pipeline []interface{}, sorter string) []interface{} {
	return append(pipeline, sorter)
}
