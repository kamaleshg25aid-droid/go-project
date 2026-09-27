package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"bloodbank"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	command := args[0]
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	file := flags.String("file", "bloodbank.json", "path to the JSON data file")
	id := flags.String("id", "", "blood unit ID")
	group := flags.String("group", "", "blood group (for collect)")
	date := flags.String("date", "", "collection date (for collect)")
	volume := flags.Int("volume", 0, "volume in milliliters")
	age := flags.Int("age", 0, "donor age (for collect)")

	if command != "collect" && command != "search" && command != "update" && command != "delete" && command != "discard" && command != "list" {
		return fmt.Errorf("unknown command %q; run with --help for usage", command)
	}
	if err := flags.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage()
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	switch command {
	case "collect":
		if *id == "" || *group == "" || *date == "" || *volume <= 0 || *age <= 0 {
			return fmt.Errorf("collect requires -id, -group, -date, -volume, and -age")
		}
	case "search", "delete", "discard":
		if *id == "" {
			return fmt.Errorf("%s requires -id", command)
		}
	case "update":
		if *id == "" || *volume <= 0 {
			return fmt.Errorf("update requires -id and a positive -volume")
		}
	}

	bank, err := bloodbank.LoadFromFile(*file)
	if err != nil {
		return fmt.Errorf("load data file: %w", err)
	}

	switch command {
	case "collect":
		result := bank.CollectUnit(bloodbank.BloodBag{
			UnitID:         *id,
			BloodGroup:     *group,
			CollectionDate: *date,
			VolumeML:       *volume,
			DonorAge:       *age,
		})
		fmt.Println(result)
		if result == "Blood bag collected" {
			return bank.SaveToFile(*file)
		}
	case "search":
		unit := bank.SearchUnit(*id)
		if unit == nil {
			fmt.Println("Blood unit not found")
		} else {
			fmt.Printf("%+v\n", *unit)
		}
	case "update":
		result := bank.UpdateVolume(*id, *volume)
		fmt.Println(result)
		if result == "Volume updated" {
			return bank.SaveToFile(*file)
		}
	case "delete", "discard":
		result := bank.DiscardUnit(*id)
		fmt.Println(result)
		if result == "Blood bag discarded" {
			return bank.SaveToFile(*file)
		}
	case "list":
		units := bank.ListUnits()
		if len(units) == 0 {
			fmt.Println("No blood units")
			return nil
		}
		for _, unit := range units {
			fmt.Printf("%+v\n", unit)
		}
	}

	return nil
}

func printUsage() {
	fmt.Println("Blood bank command line")
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/bloodbank collect -id U001 -group O+ -date 2026-09-27 -volume 450 -age 25")
	fmt.Println("  go run ./cmd/bloodbank search -id U001")
	fmt.Println("  go run ./cmd/bloodbank update -id U001 -volume 400")
	fmt.Println("  go run ./cmd/bloodbank delete -id U001")
	fmt.Println("  go run ./cmd/bloodbank list")
	fmt.Println("Optional for every command: -file path\\to\\data.json (default: bloodbank.json)")
}
