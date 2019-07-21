package engines

type Engine interface {
	UserEngine
	LinkedInEngine
	RocketChatEngine
}

type genericEngine struct {
	UserEngine
	LinkedInEngine
	RocketChatEngine
}

func NewGenericEngine(userEngine UserEngine, lEngine LinkedInEngine, rocketChatEngine RocketChatEngine) Engine {
	return &genericEngine{userEngine, lEngine, rocketChatEngine}
}
