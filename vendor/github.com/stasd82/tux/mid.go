package tux

// Circuit is a function designed to run some code before and/or after another handler.
type Circuit func(Route) Route

func warmupCircuit(ct []Circuit, route Route) Route {

	for i := len(ct) - 1; i >= 0; i-- {
		c := ct[i]
		if c != nil {
			route = c(route)
		}
	}

	return route
}
