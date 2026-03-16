package email

var Callbacks = []func(*Config, error){}

func AddCallback(cb func(*Config, error)) {
	Callbacks = append(Callbacks, cb)
}
