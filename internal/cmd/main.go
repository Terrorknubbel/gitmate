package cmd

func Main(args []string) {
	c := NewConfig()
	if len(args) == 0 {
		c.RunMenuView()
		return
	}

	err := c.execute(args)
	if err != nil {
		c.logger.Error(err)
	}
}
