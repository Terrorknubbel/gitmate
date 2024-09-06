package cmd

func Main(args []string) {
	c, err := NewConfig()

	if err != nil {
		c.logger.Error(err)
		return
	}

	if len(args) == 0 {
		c.RunMenuView()
		return
	}

	err = c.execute(args)
	if err != nil {
		c.logger.Error(err)
	}
}
