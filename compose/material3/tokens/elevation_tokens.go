package tokens

import "github.com/zodimo/go-compose/compose/ui/unit"

// import androidx.compose.ui.unit.dp

// internal object ElevationTokens {
//     val Level0 = 0.0.dp
//     val Level1 = 1.0.dp
//     val Level2 = 3.0.dp
//     val Level3 = 6.0.dp
//     val Level4 = 8.0.dp
//     val Level5 = 12.0.dp
// }

type ElevationTokenKey int

const (
	ElevationTokenKeyUnspecified ElevationTokenKey = iota
	ElevationTokenKeyLevel0
	ElevationTokenKeyLevel1
	ElevationTokenKeyLevel2
	ElevationTokenKeyLevel3
	ElevationTokenKeyLevel4
	ElevationTokenKeyLevel5
)

func (c ElevationTokenKey) String() string {
	return []string{
		"Unspecified",
		"Level0",
		"Level1",
		"Level2",
		"Level3",
		"Level4",
		"Level5",
	}[c]
}

var ElevationTokens = ElevationData{
	Level0: unit.Dp(0.0),
	Level1: unit.Dp(1.0),
	Level2: unit.Dp(3.0),
	Level3: unit.Dp(6.0),
	Level4: unit.Dp(8.0),
	Level5: unit.Dp(12.0),
}

type ElevationData struct {
	Level0 unit.Dp
	Level1 unit.Dp
	Level2 unit.Dp
	Level3 unit.Dp
	Level4 unit.Dp
	Level5 unit.Dp
}

var ElevationKeyTokens = ElevationKeyData{
	Level0: ElevationTokenKeyLevel0,
	Level1: ElevationTokenKeyLevel1,
	Level2: ElevationTokenKeyLevel2,
	Level3: ElevationTokenKeyLevel3,
	Level4: ElevationTokenKeyLevel4,
	Level5: ElevationTokenKeyLevel5,
}

type ElevationKeyData struct {
	Level0 ElevationTokenKey
	Level1 ElevationTokenKey
	Level2 ElevationTokenKey
	Level3 ElevationTokenKey
	Level4 ElevationTokenKey
	Level5 ElevationTokenKey
}
