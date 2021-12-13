/*
   Structdmaker configuration
*/

package main

import (
	"github.com/spf13/viper"
)

//All of the config vars exported
var Dbname string   //Database
var Userid string   //Database user id
var Password string //Database password

//Load the config file and extract all of the values.
func init() {
	//Setup and run Viper
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	Check(err, " Failed to read the viper configuration. ")

	//Read all of the config values
	Dbname = viper.GetString("dbname")
	Userid = viper.GetString("userid")
	Password = viper.GetString("password")
}
