package mcstatus

type JavaStatusResponse struct {
	Version            JavaVersion `json:"version"`
	Players            JavaPlayers `json:"players"`
	Description        interface{} `json:"description"`
	Favicon            string      `json:"favicon,omitempty"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat,omitempty"`
	Latency            int64       `json:"-"`
}

type JavaVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type JavaPlayers struct {
	Max    int          `json:"max"`
	Online int          `json:"online"`
	Sample []JavaPlayer `json:"sample,omitempty"`
}

type JavaPlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type BedrockStatusResponse struct {
	ServerID string `json:"serverId"`
	MOTD     string `json:"motd"`
	Protocol int    `json:"protocol"`
	Version  string `json:"version"`
	Online   int    `json:"online"`
	Max      int    `json:"max"`
	MapName  string `json:"mapName"`
	Gamemode string `json:"gamemode"`
	Latency  int64  `json:"-"`
}

func (r *JavaStatusResponse) GetLatency() int64 {
	return r.Latency
}

func (r *BedrockStatusResponse) GetLatency() int64 {
	return r.Latency
}
