package opt

import "github.com/csby/wsf/types"

var (
	optWebPath = types.Path{
		Prefix:            "/opt",
		DefaultShortenUrl: false,
		DefaultTokenType:  types.TokenTypeNone,
		DefaultTokenPlace: types.TokenPlaceQuery,
	}
	optApiPath = types.Path{
		Prefix:            "/opt.api",
		DefaultShortenUrl: false,
		DefaultTokenType:  types.TokenTypeAccountPassword,
		DefaultTokenPlace: types.TokenPlaceHeader,
		DefaultTokenUI:    types.TokenUIForAccountPassword,
	}

	webappWebPath = types.Path{
		Prefix:            "/webapp",
		DefaultShortenUrl: false,
		DefaultTokenType:  types.TokenTypeNone,
		DefaultTokenPlace: types.TokenPlaceQuery,
	}
)
