package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func init() {
	listNamespacesCmd.Flags().BoolVar(&noHeaderFlag, "no-headers", false, "Hide table headers")
	namespaceCmd.AddCommand(describeNamespace)
	namespaceCmd.AddCommand(listNamespacesCmd)
	rootCmd.AddCommand(namespaceCmd)

}

var namespaceCmd = &cobra.Command{
	Use: "namespace",
}

var listNamespacesCmd = &cobra.Command{
	Use:   "list",
	Short: "List namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)

		response, err := namespaceClient.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}
		_ = response
		if !noHeaderFlag {
			fmt.Fprintf(w, "NAME\tID\t\n")
		}

		for _, ns := range response.Namespaces {
			fmt.Fprintf(w, "%v\t%v\n", ns.GetName(), ns.GetId())
		}

		defer w.Flush()
	},
}

var describeNamespace = &cobra.Command{
	Use:     "describe",
	Short:   "Describe namespace",
	Aliases: []string{"desc"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("desc")
		w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)

		response, err := namespaceClient.GetNamespace(ctx, &nodepb.GetNamespaceRequest{Namespace: args[0]})
		if err != nil {
			fmt.Println("grpc: failed to fetch data", err)
			os.Exit(1)
		}
		fmt.Println(response.GetNamespace())
		// _ = response
		// if !noHeaderFlag {
		// 	fmt.Fprintf(w, "NAME\tID\t\n")
		// }

		// for _, ns := range response.Namespaces {
		// 	fmt.Fprintf(w, "%v\t%v\n", ns.GetName(), ns.GetId())
		// }

		defer w.Flush()
	},
}
