package bot

type GovAPI struct {
	Gov struct {
		DropGifts struct {
			ActiveGifts int `json:"activeGifts"`
			TotalGifts  int `json:"totalGifts"`
		} `json:"dropGifts"`
		Info struct {
			All int `json:"all"`
			Civ int `json:"civ"`
			Cop int `json:"cop"`
			Ems int `json:"ems"`
			Reb int `json:"reb"`
		} `json:"info"`
		Prem struct {
			Anvil        int `json:"anvil"`
			Bears        int `json:"bears"`
			Bs           int `json:"bs"`
			Cigane       int `json:"cigane"`
			Constant     int `json:"constant"`
			Cwrka        int `json:"cwrka"`
			Dl           int `json:"dl"`
			Ems          int `json:"ems"`
			Escobaro     int `json:"escobaro"`
			Flightschool int `json:"flightschool"`
			Goverment    int `json:"goverment"`
			Imperium     int `json:"imperium"`
			Judge        int `json:"judge"`
			Kasatky      int `json:"kasatky"`
			Kifo         int `json:"kifo"`
			Mgrp1        int `json:"mgrp_1"`
			Ms           int `json:"ms"`
			Narcos       int `json:"narcos"`
			Outcast      int `json:"outcast"`
			Phoenix      int `json:"phoenix"`
			Police       int `json:"police"`
			Press        int `json:"press"`
			Rimas        int `json:"rimas"`
			Sector       int `json:"sector"`
			Shark        int `json:"shark"`
			Taxi         int `json:"taxi"`
			Wolves       int `json:"wolves"`
		} `json:"prem"`
		Rule struct {
			Bank      string `json:"bank"`
			Credit    int    `json:"credit"`
			Legal     bool   `json:"legal"`
			Poor      string `json:"poor"`
			President string `json:"president"`
			Slavery   bool   `json:"slavery"`
			Tax       int    `json:"tax"`
		} `json:"rule"`
		Time int `json:"time"`
	} `json:"gov"`
	Status int `json:"status"`
}
