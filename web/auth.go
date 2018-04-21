// Copyright 2015 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package web

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailCount  = 5
	blockDuration = 24 * time.Hour
)

func getRequestIP(r *http.Request) string {
	idx := strings.LastIndex(r.RemoteAddr, ":")
	return r.RemoteAddr[:idx]
}

type basicAuthHTTPHandlerDecorator struct {
	auth        *BasicAuth
	handler     http.Handler
	handlerFunc http.HandlerFunc
	blockedIps  map[string]time.Time
	failedCount map[string]int
}

func (a *basicAuthHTTPHandlerDecorator) Unauthorized(
	w http.ResponseWriter, r *http.Request, msg string, record bool) {

	// Record failure
	if record {
		ip := getRequestIP(r)
		if _, ok := a.failedCount[ip]; !ok {
			a.failedCount[ip] = 0
		}
		if ip != "127.0.0.1" {
			// Only count for non-trusted IP.
			a.failedCount[ip]++
		}

		log.Printf("BasicAuth: IP %s failed to login, count: %d\n", ip,
			a.failedCount[ip])

		if a.failedCount[ip] >= maxFailCount {
			a.blockedIps[ip] = time.Now()
			log.Printf("BasicAuth: IP %s is blocked\n", ip)
		}
	}

	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%s", a.auth.Realm))
	http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusUnauthorized),
		msg), http.StatusUnauthorized)
}

func (a *basicAuthHTTPHandlerDecorator) IsBlocked(r *http.Request) bool {
	ip := getRequestIP(r)

	if t, ok := a.blockedIps[ip]; ok {
		if time.Now().Sub(t) < blockDuration {
			log.Printf("BasicAuth: IP %s attempted to login, blocked\n", ip)
			return true
		}
		// Unblock the user because of timeout
		delete(a.failedCount, ip)
		delete(a.blockedIps, ip)
	}
	return false
}

func (a *basicAuthHTTPHandlerDecorator) ResetFailCount(r *http.Request) {
	ip := getRequestIP(r)
	delete(a.failedCount, ip)
}

func (a *basicAuthHTTPHandlerDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.IsBlocked(r) {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusUnauthorized),
			"too many retries"), http.StatusUnauthorized)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		a.Unauthorized(w, r, "authorization failed", false)
		return
	}

	pass, err := a.auth.Authenticate(username, password)
	if !pass {
		a.Unauthorized(w, r, err.Error(), true)
		return
	}
	a.ResetFailCount(r)

	if a.handler != nil {
		a.handler.ServeHTTP(w, r)
	} else {
		a.handlerFunc(w, r)
	}
}

// BasicAuth is a class that provide  WrapHandler and WrapHandlerFunc, which
// turns a http.Handler to a HTTP basic-auth enabled http handler.
type BasicAuth struct {
	Realm   string
	secrets map[string]string
}

// NewBasicAuth creates a BasicAuth object
func NewBasicAuth(realm, htpasswd string) *BasicAuth {
	secrets := make(map[string]string)

	f, err := os.Open(htpasswd)
	if err != nil {
		return &BasicAuth{realm, secrets}
	}

	b := bufio.NewReader(f)
	for {
		line, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}
		if line[0] == '#' {
			continue
		}
		parts := strings.Split(string(line), ":")
		if len(parts) != 2 {
			continue
		}
		matched, err := regexp.Match("^\\$2[ay]\\$.*$", []byte(parts[1]))
		if err != nil {
			panic(err)
		}
		if !matched {
			log.Printf("BasicAuth: user %s: password encryption scheme "+
				"not supported, ignored.\n", parts[0])
			continue
		}
		secrets[parts[0]] = parts[1]
	}

	return &BasicAuth{realm, secrets}
}

// WrapHandler wraps an http.Hanlder and provide HTTP basic-a.
func (a *BasicAuth) WrapHandler(h http.Handler) http.Handler {
	return &basicAuthHTTPHandlerDecorator{a, h, nil,
		make(map[string]time.Time), make(map[string]int)}
}

// WrapHandlerFunc wraps an http.HanlderFunc and provide HTTP basic-a.
func (a *BasicAuth) WrapHandlerFunc(h http.HandlerFunc) http.Handler {
	return &basicAuthHTTPHandlerDecorator{a, nil, h,
		make(map[string]time.Time), make(map[string]int)}
}

// Authenticate authenticate an user with the provided user and passwd.
func (a *BasicAuth) Authenticate(user, passwd string) (bool, error) {
	deniedError := errors.New("permission denied")

	passwdHash, ok := a.secrets[user]
	if !ok {
		return false, deniedError
	}

	if bcrypt.CompareHashAndPassword([]byte(passwdHash), []byte(passwd)) != nil {
		return false, deniedError
	}

	return true, nil
}
