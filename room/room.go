package room

type roomOrchestrator struct {
	repository
}

func Room(repository repository) roomOrchestrator {
	return roomOrchestrator{
		repository,
	}
}

func (r roomOrchestrator) Enter() error {
	return nil
}

func (r roomOrchestrator) Exit() error {
	return nil
}
