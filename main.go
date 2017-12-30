package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/genevievelesperance/az-automation/az"
	flags "github.com/jessevdk/go-flags"
)

type args struct {
	Account              string `required:"true" short:"a" long:"account"                description:"Your account id or name. Use 'az account list' to see your accounts."`
	DisplayName          string `required:"true" short:"d" long:"display-name"           description:"Display name for application. Must be unique."`
	IdentifierUri        string `required:"true" short:"i" long:"identifier-uri"         description:"Must be unique."`
	CredentialOutputFile string `required:"true" short:"c" long:"credential-output-file" description:"Must be unique."                                                      default:"creds"`
}

func main() {
	log.SetFlags(0)

	var a args
	parser := flags.NewParser(&a, flags.HelpFlag|flags.PrintErrors)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(0)
	}

	path, err := exec.LookPath("az")
	if err != nil {
		log.Fatalf("Failed to find az (azure-cli): %s", err)
	}

	binary := az.NewCLI(path)

	cli := az.NewAz(binary, a.Account, a.DisplayName, a.IdentifierUri, a.CredentialOutputFile)

	err = cli.ValidVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Checked version of azure-cli.")

	err = cli.LoggedIn()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Checked you are logged in to azure-cli (`az`).")

	err = cli.GetSubscriptionAndTenantId()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Retrieved subscription and tenant id.")

	err = cli.AppExists()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Confirmed application name is not already taken.")

	err = cli.CreateApplication()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created application.")

	err = cli.CreateServicePrincipal()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created service principal.")

	err = cli.AssignContributorRole()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Assigned contributor role to service principal.")

	err = cli.WriteCredentials()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Wrote credentials to output file.")
}
