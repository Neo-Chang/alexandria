/*
 * Alexandria CMDB - Open source common.management database
 * Copyright (C) 2014  Ryan Armstrong <ryan@cavaliercoder.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * package controllers
 */
package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cavaliercoder/alexandria/models"
	"github.com/cavaliercoder/alexandria/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type DatabaseController struct {
	controller
}

func (c *DatabaseController) GetPath() string {
	return "/databases"
}

func (c *DatabaseController) InitRoutes(r martini.Router) error {
	r.Get("/", c.getDatabases)
	r.Post("/", binding.Bind(models.Database{}), c.createDatabase)
	r.Get("/:shortname", c.getDatabaseByShortName)

	return nil
}

func (c *DatabaseController) getDatabases(r *services.ApiContext) {
	r.Render(http.StatusOK, r.AuthTenant.Databases)
}

func (c *DatabaseController) getDatabaseByShortName(r *services.ApiContext, params martini.Params) {
	for _, db := range r.AuthTenant.Databases {
		if db.ShortName == params["shortname"] {
			r.Render(http.StatusOK, db)
			return
		}
	}

	r.NotFound()
}

func (c *DatabaseController) createDatabase(database models.Database, r *services.ApiContext) {
	database.Init(r.DB.NewId())
	database.TenantId = r.AuthUser.TenantId

	// TODO: Append the database to a tenant instead of a collection
	err := r.DB.
	
	// Create database entry
	err := r.DB.Insert("databases", &database)
	if err != nil {
		log.Panic(err)
	}

	// Create actual database
	err = r.DB.CreateDatabase(database.Backend)

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/databases/%s", database.ShortName))
	r.Render(http.StatusCreated, "")
}