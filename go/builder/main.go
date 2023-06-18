package main

func main() {
	// K8s Builder pattern 學習
}

type K8sFake struct {
	condition string
	eventType string
	recorder  string
	configs   configs
}

func Create(condition, eventType, recorder string, os ...Option) error {
	opts := &configs{
		force:     false,
		onSuccess: func() {},
		onError:   nil,
	}

	// apply all the optional configs
	for _, o := range os {
		o(opts)
	}
	// check required fields

	// update conditions here

	// handle error here
	if opts.err != nil {
		return opts.onError()
	}

	// eveutally, call success function
	opts.onSuccess()

	return nil
}

// configs 可選的配置
type configs struct {
	force     bool         // is Force update
	onError   func() error // err handler
	onSuccess func()       // success handler

	err error // error
}

type Option func(*configs)

func ForceUpdate(force bool) Option { return func(c *configs) { c.force = force } }

func OnErr(onErr func() error) Option { return func(c *configs) { c.onError = onErr } }

func OnSuccess(onSucc func()) Option { return func(c *configs) { c.onSuccess = onSucc } }
