package main

import "net/http"

func (app *application) usersSignupPost(w http.ResponseWriter, r *http.Request) {}

func (app *application) usersLoginPost(w http.ResponseWriter, r *http.Request) {}

func (app *application) usersLogoutPost(w http.ResponseWriter, r *http.Request) {}

func (app *application) usersMeGet(w http.ResponseWriter, r *http.Request) {}
