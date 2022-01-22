package config

type Office struct {
	LOG LOG
	GW  GW
}

type GW struct {
	IP   string `default:"gw.globalzt.com"`
	PORT string `default:"31580"`
}

func (c *Office) Load() {
	env.ignorePrefix = true
	env.Fill(c)
}
