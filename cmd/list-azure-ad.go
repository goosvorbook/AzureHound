// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listAzureADCmd)
}

var listAzureADCmd = &cobra.Command{
	Use:               "az-ad",
	Long:              "Lists All Azure AD Entities",
	PersistentPreRunE: persistentPreRunE,
	Run:               listAzureADCmdImpl,
	SilenceUsage:      true,
}

func listAzureADCmdImpl(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		exit(fmt.Errorf("unsupported subcommand: %v", args))
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure ad objects...")
	start := time.Now()
	stream := listAllAD(ctx, azClient)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listAllAD(ctx context.Context, client client.AzureClient) <-chan any {
	var (
		devices  = make(chan any)
		devices2 = make(chan any)

		groups  = make(chan any)
		groups2 = make(chan any)
		groups3 = make(chan any)

		roles  = make(chan any)
		roles2 = make(chan any)

		servicePrincipals  = make(chan any)
		servicePrincipals2 = make(chan any)
		servicePrincipals3 = make(chan any)

		tenants = make(chan any)
	)

	// Enumerate Apps, AppOwners and AppMembers
	appChans := pipeline.TeeFixed(ctx.Done(), listApps(ctx, client), 2)
	apps := pipeline.ToAny(ctx.Done(), appChans[0])
	appOwners := pipeline.ToAny(ctx.Done(), listAppOwners(ctx, client, appChans[1]))

	// Enumerate Devices and DeviceOwners
	pipeline.Tee(ctx.Done(), listDevices(ctx, client), devices, devices2)
	deviceOwners := listDeviceOwners(ctx, client, devices2)

	// Enumerate Groups, GroupOwners and GroupMembers
	pipeline.Tee(ctx.Done(), listGroups(ctx, client), groups, groups2, groups3)
	groupOwners := listGroupOwners(ctx, client, groups2)
	groupMembers := listGroupMembers(ctx, client, groups3)

	// Enumerate ServicePrincipals and ServicePrincipalOwners
	pipeline.Tee(ctx.Done(), listServicePrincipals(ctx, client), servicePrincipals, servicePrincipals2, servicePrincipals3)
	servicePrincipalOwners := listServicePrincipalOwners(ctx, client, servicePrincipals2)

	// Enumerate Tenants
	pipeline.Tee(ctx.Done(), listTenants(ctx, client), tenants)

	// Enumerate Users
	users := listUsers(ctx, client)

	// Enumerate Roles and RoleAssignments
	pipeline.Tee(ctx.Done(), listRoles(ctx, client), roles, roles2)
	roleAssignments := listRoleAssignments(ctx, client, roles2)

	// Enumerate AppRoleAssignments
	appRoleAssignments := listAppRoleAssignments(ctx, client, servicePrincipals3)

	return pipeline.Mux(ctx.Done(),
		appOwners,
		appRoleAssignments,
		apps,
		deviceOwners,
		devices,
		groupMembers,
		groupOwners,
		groups,
		roleAssignments,
		roles,
		servicePrincipalOwners,
		servicePrincipals,
		tenants,
		users,
	)
}
