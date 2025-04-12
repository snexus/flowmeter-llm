/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"github.com/spf13/cobra"
	"github.com/snexus/wmeter/entities"
	"github.com/snexus/wmeter/usecases"
	"github.com/snexus/wmeter/llm"
	"github.com/snexus/wmeter/db"
)

const prompt string = "What value is displayed in the water meter at the center (black and red)? The value contains 8 digits. The image can be upside down or rotated."


var modelName string
var openaiEndpoint string
var nImages int

// ingestCmd represents the ingest command
var ingestCmd = &cobra.Command{
	Use:   "ingest <folder>",
	Short: "Ingests and parses water meter readings and saves them to the database.",
	Long: `This command scans a folder for images of water meter and uses LLM to extract the readings from the images.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning for new images and ingesting into the database.")

	imagesToProcess := []entities.ImageMetadata{}
	fsImages := usecases.FetchLatestImageFilesFromDisk(args[0], "*.jpg", nImages)
	// fmt.Println(fsImages)

	db := db.InitDB("./app.db")
	defer db.Db.Close()

	dbImages, err := db.QueryLatestDays(nImages*2)

	imageMap := make(map[string]entities.ImageMetadata)
	for _, image := range dbImages {
		imageMap[image.Hash] = image
	}
	// fmt.Print("dbImages: ", imageMap)

	if err != nil {
		log.Fatal(err)
	}

	for hash, fsImage := range fsImages {
		_, ok := imageMap[hash]
		if !ok {
			imagesToProcess = append(imagesToProcess, fsImage)
		}
	}

	fmt.Println("\nNew imagesToProcess: ", imagesToProcess)

	if len(imagesToProcess) > 0 {

		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY environment variable is not set.")
			os.Exit(1)
		}

		client := llm.GetOpenAIClient(apiKey, openaiEndpoint)
		analyzedImages := usecases.AnalyzeMultipleImages(client, imagesToProcess, prompt, modelName)
		for _, image := range analyzedImages {
			fmt.Println("Recording image: ", image)
			err := db.InsertRecord(image)
			if err != nil {
				fmt.Printf("Error inserting record into the database: %v\n", err)
				continue
			}
		}
	} else {
		fmt.Println("Nothing new to process. Exiting.")
	}

	},
}

func init() {
	rootCmd.AddCommand(ingestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ingestCmd.PersistentFlags().String("foo", "", "A help for foo")
	ingestCmd.Flags().StringVarP(&modelName, "model", "m", "gemini-2.0-flash", "Model name to use for parsing the images.")
	ingestCmd.Flags().StringVarP(&openaiEndpoint, "endpoint", "e", "", "OpenAI endpoint to use.")
	_ = ingestCmd.MarkFlagRequired("endpoint")
	ingestCmd.Flags().IntVarP(&nImages, "nimg", "n", 6 , "Number of most recent images to process.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ingestCmd.Flags().String("model", "t", false, "Help message for toggle")
}
