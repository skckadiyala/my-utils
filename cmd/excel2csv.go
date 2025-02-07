/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/skckadiyala/my-utils/pkg/utils"
	"github.com/spf13/cobra"
)

// excel2csvCmd represents the excel2csv command
var xlsheet2csvCmd = &cobra.Command{
	Use:   "xlsheet2csv",
	Short: "Convert an excel sheet to csv",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.ExcelSheet2CSV(xlsxFile, sheetName)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

var excel2csvCmd = &cobra.Command{
	Use:   "excel2csv",
	Short: "Convert all the sheets to individual csv's",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.Excel2CSV(xlsxFile)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

var csv2excelCmd = &cobra.Command{
	Use:   "csv2excel",
	Short: "Convert the csv file to excel worksheet",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.CSV2Excel(csvfile, xlsxFile)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

var csv2jsonCmd = &cobra.Command{
	Use:   "csv2json",
	Short: "Convert the csv file to Json file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := utils.CSV2Json(csvfile)
		return err
		// return utils.CSV2Json(csvfile, jsonfile)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

var push2splunkCmd = &cobra.Command{
	Use:   "provar2splunk",
	Short: "Push Provar results from Excel to Splunk",
	Long: `Push Provar Results to Splunk, In this the excel results sheet is converted to json 
	and push the json results to splunk. For example:

	myutils convert provar2splunk -H <splunkhost> -u <splunk user> -p <splunk password>  -s <source splunk> -x resources/TestExecutionReport.xlsx

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.ProvarResults2Splunk(xlsxFile, splunkHost, splunkPort, userName, password, source, index)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

var jMeter2SplunkCmd = &cobra.Command{
	Use:   "jMeter2SplunkCmd",
	Short: "Push Provar results from JMeter csv to Splunk",
	Long: `Push Provar Results to Splunk, In this the excel results sheet is converted to json 
	and push the json results to splunk. For example:

	myutils convert provar2splunk -H <splunkhost> -u <splunk user> -p <splunk password>  -s <source splunk> -x resources/TestExecutionReport.xlsx

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.JMeterResults2Splunk(csvfile, splunkHost, splunkPort, userName, password, source, index)
		// return xlsx2csv()
		// fmt.Println("excel2csv called")

	},
}

func init() {
	convertCmd.AddCommand(xlsheet2csvCmd)
	convertCmd.AddCommand(excel2csvCmd)
	convertCmd.AddCommand(csv2excelCmd)
	convertCmd.AddCommand(csv2jsonCmd)
	convertCmd.AddCommand(push2splunkCmd)
	convertCmd.AddCommand(jMeter2SplunkCmd)

	xlsheet2csvCmd.Flags().StringVarP(&xlsxFile, "xlsxfile", "x", "", "The filename of the excel file")
	xlsheet2csvCmd.MarkFlagRequired("xlsxfile")
	// xlsheet2csvCmd.Flags().StringVarP(&csvfile, "csvfile", "c", "", "The filename of the csv file")
	// xlsheet2csvCmd.MarkFlagRequired("csvfile")
	xlsheet2csvCmd.Flags().StringVarP(&sheetName, "sheetName", "s", "", "The name of the sheet in excel to convert")
	xlsheet2csvCmd.MarkFlagRequired("sheetName")

	excel2csvCmd.Flags().StringVarP(&xlsxFile, "xlsxfile", "x", "", "The filename of the excel file")
	excel2csvCmd.MarkFlagRequired("xlsxfile")

	csv2excelCmd.Flags().StringVarP(&xlsxFile, "xlsxfile", "x", "", "The filename of the excel file")
	csv2excelCmd.MarkFlagRequired("xlsxfile")
	csv2excelCmd.Flags().StringVarP(&csvfile, "csvfile", "c", "", "The filename of the csv file")
	csv2excelCmd.MarkFlagRequired("csvfile")

	csv2jsonCmd.Flags().StringVarP(&csvfile, "csvfile", "c", "", "The filename of the csv file")
	csv2jsonCmd.MarkFlagRequired("csvfile")
	// csv2jsonCmd.Flags().StringVarP(&jsonfile, "jsonfile", "j", "", "The filename of the Json file")
	// csv2jsonCmd.MarkFlagRequired("jsonfile")

	push2splunkCmd.Flags().StringVarP(&xlsxFile, "xlsxFile", "x", "", "Provide the filename of the excel file")
	push2splunkCmd.MarkFlagRequired("jsonfile")
	push2splunkCmd.Flags().StringVarP(&splunkHost, "splunkHost", "H", "", "Provide the host name")
	push2splunkCmd.MarkFlagRequired("splunkHost")
	push2splunkCmd.Flags().StringVarP(&splunkPort, "splunkPort", "P", "8089", "Provide the splunk port name")
	// push2splunkCmd.MarkFlagRequired("splunkPort")
	push2splunkCmd.Flags().StringVarP(&userName, "userName", "u", "", "Provide the username of splunk")
	push2splunkCmd.MarkFlagRequired("userName")
	push2splunkCmd.Flags().StringVarP(&password, "password", "p", "", "Provide the password of splunk")
	push2splunkCmd.MarkFlagRequired("password")
	push2splunkCmd.Flags().StringVarP(&source, "source", "s", "", "Provide the splunk source name")
	push2splunkCmd.MarkFlagRequired("source")
	push2splunkCmd.Flags().StringVarP(&index, "index", "i", "sdlc", "Provide the splunk index file")

	jMeter2SplunkCmd.Flags().StringVarP(&csvfile, "csvFile", "c", "", "Provide the filename of the Jmeter Results csv file")
	jMeter2SplunkCmd.MarkFlagRequired("jsonfile")
	jMeter2SplunkCmd.Flags().StringVarP(&splunkHost, "splunkHost", "H", "", "Provide the host name")
	jMeter2SplunkCmd.MarkFlagRequired("splunkHost")
	jMeter2SplunkCmd.Flags().StringVarP(&splunkPort, "splunkPort", "P", "8089", "Provide the splunk port name")
	// push2splunkCmd.MarkFlagRequired("splunkPort")
	jMeter2SplunkCmd.Flags().StringVarP(&userName, "userName", "u", "", "Provide the username of splunk")
	jMeter2SplunkCmd.MarkFlagRequired("userName")
	jMeter2SplunkCmd.Flags().StringVarP(&password, "password", "p", "", "Provide the password of splunk")
	jMeter2SplunkCmd.MarkFlagRequired("password")
	jMeter2SplunkCmd.Flags().StringVarP(&source, "source", "s", "", "Provide the splunk source name")
	jMeter2SplunkCmd.MarkFlagRequired("source")
	jMeter2SplunkCmd.Flags().StringVarP(&index, "index", "i", "sdlc", "Provide the splunk index file")

	// push2splunkCmd.MarkFlagRequired("index")

	// -H splnkfrprdvh3 -u SplunkApi -p s#uVUheSwA5W -s CDW_ECommerceWeb_TestAutomation -i sdlc
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// excel2csvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// excel2csvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
