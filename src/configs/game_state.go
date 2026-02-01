package configs

type GameStateEnum int

const (
	GameStateInit GameStateEnum = iota
	GameStatePlaying
	GameStatePaused
	GameStateGameOver
)

type GameSceneEnum int

const (
	GameSceneLoadingScreen GameSceneEnum = iota
	GameSceneMainMenu
	GameScenePlay
	GameSceneOver
	GameScenePaused
	GameSceneSettings
	GameSceneCredits
	GameSceneHelp
	GameSceneHighScores
	GameSceneAchievements
)

type GameEventEnum int

const (
	GameEventMaskConsumed GameEventEnum = iota
)

type GameStateManager struct {
	state        GameStateEnum
	events       map[GameEventEnum][]func()
	SurvivalTime float64
	RowsCleared  int
}

var globalGameState = &GameStateManager{
	state:  GameStatePlaying,
	events: make(map[GameEventEnum][]func()),
}

func GameState() *GameStateManager {
	return globalGameState
}

func (g *GameStateManager) Update() {
	if g != globalGameState {
		return
	}
}

func (g *GameStateManager) GetGameState() GameStateEnum {
	return g.state
}

func (g *GameStateManager) SetGameState(state GameStateEnum) {
	g.state = state
}

func (g *GameStateManager) TriggerEvent(eventName GameEventEnum) {
	if event, ok := g.events[eventName]; ok {
		for _, callback := range event {
			callback()
		}
	}
}

func (g *GameStateManager) RegisterEventHandler(eventName GameEventEnum, callback func()) {
	if _, ok := g.events[eventName]; !ok {
		g.events[eventName] = make([]func(), 0)
	}
	g.events[eventName] = append(g.events[eventName], callback)
}

func (g *GameStateManager) AddSurvivalTime(delta float64) {
	if g.GetGameState() == GameStatePlaying {
		g.SurvivalTime += delta
	}
}

func (g *GameStateManager) IncrementRowsCleared() {
	g.RowsCleared++
}

func (g *GameStateManager) ResetScore() {
	g.SurvivalTime = 0
	g.RowsCleared = 0
}
