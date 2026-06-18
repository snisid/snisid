package domain

type WeaponsByGang struct {
	GangID      string `json:"gang_id"`
	GangName    string `json:"gang_name"`
	WeaponCount int    `json:"weapon_count"`
}

type WeaponsByOrigin struct {
	OriginCountry string `json:"origin_country"`
	WeaponCount   int    `json:"weapon_count"`
}

type TraffickingRoute struct {
	OriginCountry    string   `json:"origin_country"`
	TransitCountries []string `json:"transit_countries"`
	ImportMethod     string   `json:"import_method"`
	WeaponCount      int      `json:"weapon_count"`
}
