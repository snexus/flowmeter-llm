/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"github.com/snexus/wmeter/usecases"
	"github.com/spf13/cobra"
	"github.com/snexus/wmeter/db"
)

var nDays int

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze the water meter readings and calculate the average daily consumption.",
	Long: `This command analyzes the water meter readings and calculates the average daily consumption based on the readings and the time between the readings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Analyzing flow meter readings...")

	db := db.InitDB("./app.db")
	defer db.Db.Close()

	imageData, err := db.QueryLatestDays(nDays)

	if err != nil {
		log.Fatal("Can't get data from the database.")
	}
	usecases.CalculateFlow(imageData, nDays)
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().IntVarP(&nDays, "ndays", "n", 10 , "Number of last days to analyze")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyzeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyzeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
