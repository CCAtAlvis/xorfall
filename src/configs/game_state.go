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

type GameStateManager struct {
	state GameStateEnum
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

	switch g.state {
	case GameStatePlaying:
		// TODO: Implement playing state
	}
}

func (g *GameStateManager) GetGameState() GameStateEnum {
	return g.state
}

func (g *GameStateManager) SetGameState(state GameStateEnum) {
	g.state = state
}
