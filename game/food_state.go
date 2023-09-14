package game

type FoodState interface {
	UpdateState(input *foodStateInput)
	OnUpdateState(fn func())
	GetPosition() foodPosition
}

type foodPosition struct {
	X int
	Y int
}

type foodState struct {
	position        foodPosition
	onUpdateHandler func()
}

type foodStateInput struct {
	position *foodPosition
}

func NewFoodState() FoodState {
	return &foodState{}
}

func (fs *foodState) UpdateState(input *foodStateInput) {
	if input.position != nil {
		fs.position = *input.position
	}

	fs.dispatchUpdateEvent()
}

func (fs *foodState) dispatchUpdateEvent() {
	if fs.onUpdateHandler != nil {
		fs.onUpdateHandler()
	}
}

func (fs *foodState) OnUpdateState(fn func()) {
	fs.onUpdateHandler = fn
}

func (fs *foodState) GetPosition() foodPosition {
	return fs.position
}
