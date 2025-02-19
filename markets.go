package spotifind

import "slices"

// This file contains the list of markets that Spotify supports.
// All the markets are defined as ISO 3166-1 alpha-2 country codes.

var (
	marketsPopular = []string{
		"US", "GB", "DE", "FR", "CA", "AU", "JP", "IT", "NL", "ES", "BR", "IN", "KR",
	}

	marketsAsia = []string{
		"AM", "BH", "BD", "BT", "BN", "KH", "GE", "HK", "ID", "IQ", "IL", "JO", "KW", "KG", "LA", "LB", "MO", "MY", "MV", "MN", "NP", "OM", "PK", "PS", "PH", "QA", "SA", "SG", "LK", "TW", "TJ", "TH", "TL", "AE", "UZ", "VN",
	}
	marketsEurope = []string{
		"AL", "AD", "AT", "PT", "BY", "BE", "BA", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "GR", "HU", "IS", "IE", "KZ", "XK", "LV", "LI", "LT", "LU", "PT", "MT", "MD", "MC", "ME", "MK", "NO", "PL", "PT", "RO", "SM", "RS", "SK", "SI", "SE", "CH", "TR", "UA",
	}
	marketsNorthAmerica = []string{
		"AG", "BS", "BB", "BZ", "CR", "CW", "DM", "DO", "SV", "GD", "GT", "HT", "HN", "JM", "MX", "NI", "PA", "PR", "KN", "LC", "VC", "TT",
	}
	marketsSouthAmerica = []string{
		"AR", "BO", "CL", "CO", "EC", "GY", "PY", "PE", "SR", "UY", "VE",
	}
	marketsOceania = []string{
		"FJ", "KI", "MH", "FM", "NR", "NZ", "PW", "PG", "WS", "SB", "TO", "TV", "VU",
	}
	marketsAfrica = []string{
		"DZ", "AO", "BJ", "BW", "BF", "BI", "CM", "CV", "TD", "KM", "CI", "CD", "DJ", "EG", "ET", "GQ", "SZ", "GA", "GM", "GH", "GN", "GW", "KE", "LS", "LR", "LY", "MG", "MW", "ML", "MR", "MU", "MA", "MZ", "NA", "NE", "NG", "CG", "RW", "ST", "SN", "SC", "SL", "ZA", "TZ", "TG", "TN", "UG", "ZM", "ZW",
	}

	marketsUnpopular = slices.Concat(marketsEurope, marketsAsia, marketsOceania, marketsNorthAmerica, marketsSouthAmerica, marketsAfrica)
	MarketsAll       = slices.Concat(marketsPopular, marketsUnpopular)
)
