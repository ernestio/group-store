/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
	"time"

	"github.com/nats-io/nats"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHandler(t *testing.T) {
	setupNats()
	n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","groups","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	setupPg()
	startHandler()

	Convey("Scenario: getting a group", t, func() {
		setupTestSuite()
		Convey("Given the group does not exist on the database", func() {
			msg, err := n.Request("group.get", []byte(`{"id":"32"}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldEqual, nil)
		})
	})

	Convey("Scenario: deleting a group", t, func() {
		setupTestSuite()
	})

	Convey("Scenario: group set", t, func() {
		setupTestSuite()
	})

	Convey("Scenario: find datacenters", t, func() {
		setupTestSuite()
	})
}
