package constants

const (
	ResourceKeyUserID     = "user-id"
	ResourceKeyPlaylistID = "playlist-id"
	ResourceKeyGameID     = "game-id"
	ResourceKeyUsername   = "username"
	ResourceKeyPostID     = "post-id"
)

const (
	RequestJSON = "json"
	RequestData = "data"
)

type PublicResponse struct {
	Msg    *string `json:"message"`
	Status int     `json:"status"`
}

type FilterGroup struct {
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Extreme bool     `json:"extreme"`
}

type ResponseFilterGroups struct {
	Groups []*FilterGroup `json:"groups"`
}

var FilterGroupSeizureWarning = &FilterGroup{
	Name:    "Seizure Warning",
	Tags:    []string{"Seizure Warning"},
	Extreme: false,
}

var FilterGroupPornography = &FilterGroup{
	Name: "Pornography",
	Tags: []string{"Anal",
		"Anal Insertion",
		"BDSM",
		"Cartoon Porn",
		"Adult",
		"Anilingus",
		"Fingering",
		"Incest",
		"Oral",
		"Sexual Content",
		"Cunnilingus",
		"Fellatio",
		"Footjob",
		"Handjob",
		"Hypnosis",
		"Infantilism",
		"Inflation",
		"Interspecies",
		"Masturbation",
		"Paizuri",
		"Pregnancy",
		"Sex Toys",
		"Spanking",
		"Tentacles",
		"Touching",
		"Tribadism",
		"Urination",
		"Vaginal",
		"Vaginal Insertion",
		"Futanari",
		"Male Futanari",
		"Gynomorph",
		"Andromorph",
		"Oviposition",
		"Intersex",
		"Breast Milking",
		"Porn",
		"Hentai",
		"Group",
		"Solo",
		"Cannibalism",
		"Enema",
		"Frottage",
		"Kabeshiri",
		"Macrophilia",
		"Obesity",
		"Podophilia",
		"Quicksand",
		"Tickling",
		"Weight Gain",
		"Gloryhole",
		"Multiple Penises",
		"Ambiguous Penetration",
		"Self Oral"},
	Extreme: true,
}

var FilterGroupViolence = &FilterGroup{
	Name: "Violence",
	Tags: []string{"Gore",
		"Strong Violence",
		"Strong Language"},
	Extreme: true,
}

var FilterGroupBigotry = &FilterGroup{
	Name: "Bigotry",
	Tags: []string{"Homophobia",
		"Stereotyping",
		"Racism",
		"Transphobia"},
	Extreme: true,
}

var FilterGroupPornographyExtreme = &FilterGroup{
	Name: "Pornography (Extreme)",
	Tags: []string{"Bestiality",
		"Cannibalism",
		"Enema",
		"Fisting",
		"Flatulence",
		"Necrophilia",
		"Scat",
		"Vomit",
		"Vore",
		"Sexual Violence"},
	Extreme: true,
}

var FilterGroupMatureTopics = &FilterGroup{
	Name: "Otherwise Mature Topics",
	Tags: []string{"Drugs",
		"Reproductive Health",
		"Addiction",
		"Heavy Themes",
		"Suicide",
		"Nudity",
		"Moderate Language",
		"Sexual Harassment"},
	Extreme: true,
}

var FilterGroups = []*FilterGroup{
	FilterGroupSeizureWarning,
	FilterGroupPornography,
	FilterGroupViolence,
	FilterGroupBigotry,
	FilterGroupPornographyExtreme,
	FilterGroupMatureTopics,
}

func GetFilterGroups() []*FilterGroup {
	return FilterGroups
}
