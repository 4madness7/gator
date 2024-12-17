package config

func (c *Config) SetUser(currentUser string) error{
    c.CurrentUserName = currentUser
    err := write(*c)
    if err != nil {
        return err
    }
    return nil
}
