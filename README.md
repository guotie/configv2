# Config Readme

------

This package is used to get config params in json format text file.

------

## Usage

```golang

// AppConfig app config
type AppConfig struct {
	Dbhost     string `json:"dbhost" default:"localhost"`
	Dbport     string `json:"dbport" default:"3306"`
	Dbuser     string `json:"dbuser"`
	Dbpassword string `json:"dbpassword"`
	Dbname     string `json:"dbname"`
	Dbtyp      string `json:"dbtype" default:"mysql"`
	BindAddr   string `json:"bindAddr" default:":80"`
	SecretKey  string `json:"secretkey"`
}

var app AppConfig

// fn is config.json path
func readConfig(fn string) error {
	c := configv2.NewFileConfig(fn)

	err := c.Read(&app)
	if err != nil {
		return err
	}

	return nil
}

```
