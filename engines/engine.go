package engines

type Engine interface {
	UserEngine
	LinkedInEngine
}

type genericEngine struct {
	UserEngine
	LinkedInEngine
}

func NewGenericEngine(userEngine UserEngine, lEngine LinkedInEngine) Engine {
	return &genericEngine{userEngine, lEngine}
}
