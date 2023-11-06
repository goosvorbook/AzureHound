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

package azure

import "github.com/gofrs/uuid"

// Represents an application role that can be requested by (and granted to) a client application, or that can be used to
// assign an application to users or groups in a specified role.
//
// An app role assignment is a relationship between the assigned principal (a user, a group, or a service principal),
// a resource application (the app's service principal) and an app role defined on the resource application.
//
// With appRoleAssignments, app roles can be assigned to users, groups, or other applications' service principals.
// For more detail see https://docs.microsoft.com/en-us/graph/api/resources/approleassignment?view=graph-rest-1.0
type AppRoleAssignment struct {
	AppRoleId            uuid.UUID `json:"appRoleId"`
	CreatedDateTime      string    `json:"createdDateTime"`
	Id                   string    `json:"id"`
	PrincipalDisplayName string    `json:"principalDisplayName"`
	PrincipalId          uuid.UUID `json:"principalId"`
	PrincipalType        string    `json:"principalType"`
	ResourceDisplayName  string    `json:"resourceDisplayName"`
	ResourceId           string    `json:"resourceId"`
}

type AppRoleAssignmentList struct {
	Count    int                 `json:"@odata.count,omitempty"`    // The total count of all results
	NextLink string              `json:"@odata.nextLink,omitempty"` // The URL to use for getting the next set of values.
	Context  string              `json:"@odata.context,omitempty"`
	Value    []AppRoleAssignment `json:"value"` // A list of role assignments.
}

type AppRoleAssignmentResult struct {
	Error error
	Ok    AppRoleAssignment
}
