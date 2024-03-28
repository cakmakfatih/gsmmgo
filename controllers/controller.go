package controllers

type Controller interface {
	register()
}

func RegisterControllers(controllers []Controller) {
	for _, c := range controllers {
		c.register()
	}
}
