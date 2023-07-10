package utils

type ExceptionStruct struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}
type Exception interface{}

func Throw(up Exception) {
	panic(up)
}
func (this ExceptionStruct) Do() {
	if this.Finally != nil {

		defer this.Finally()
	}
	if this.Catch != nil {
		defer func() {
			if e := recover(); e != nil {
				this.Catch(e)
			}
		}()
	}
	this.Try()
}
