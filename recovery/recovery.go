package recovery

import (
	"github.com/ExtraWhy/internal-libs/models/games"
	"github.com/ExtraWhy/internal-libs/models/user"
)

var sGames = map[user.User]games.Games{}
