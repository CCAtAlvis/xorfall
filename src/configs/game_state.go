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
	state  GameStateEnum
	events map[GameEventEnum][]func()
}

var globalGameState = &GameStateManager{
	state: GameStatePlaying,
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
