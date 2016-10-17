package bulldozer

type Task interface {
	Run(interface{}) interface{}
}
