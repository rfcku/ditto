package utils

type ToView interface {
	View() error
}

func ToViewObject(v ToView) error {
	return v.View()
}
