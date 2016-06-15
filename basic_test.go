/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
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

		Convey("Given the group exists on the database", func() {
			createEntities(1)
			e := Entity{}
			db.First(&e)
			id := fmt.Sprint(e.ID)

			msg, err := n.Request("group.get", []byte(`{"id":`+id+`}`), time.Second)
			output := Entity{}
			json.Unmarshal(msg.Data, &output)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(err, ShouldEqual, nil)
		})

		Convey("Given the group exists on the database and searching by name", func() {
			createEntities(1)
			e := Entity{}
			db.First(&e)

			msg, err := n.Request("group.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			json.Unmarshal(msg.Data, &output)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("Scenario: deleting a group", t, func() {
		setupTestSuite()
		Convey("Given the group does not exist on the database", func() {
			msg, err := n.Request("group.del", []byte(`{"id":32}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldEqual, nil)
		})

		Convey("Given the group exists on the database", func() {
			createEntities(1)
			last := Entity{}
			db.First(&last)
			id := fmt.Sprint(last.ID)

			msg, err := n.Request("group.del", []byte(`{"id":`+id+`}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.DeletedMessage))
			So(err, ShouldEqual, nil)

			deleted := Entity{}
			db.First(&deleted, id)
			So(deleted.ID, ShouldEqual, 0)
		})
	})

	Convey("Scenario: group set", t, func() {
		setupTestSuite()
		Convey("Given we don't provide any id as part of the body", func() {
			Convey("Then it should return the created record and it should be stored on DB", func() {
				msg, err := n.Request("group.set", []byte(`{"name":"foo"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "foo")
				So(err, ShouldEqual, nil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "foo")
			})
		})

		Convey("Given we provide an unexisting id", func() {
			Convey("Then we should receive a not found message", func() {
				msg, err := n.Request("group.set", []byte(`{"id": 1000, "name":"foo"}`), time.Second)
				So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Given we provide an existing id", func() {
			createEntities(1)
			e := Entity{}
			db.First(&e)
			id := fmt.Sprint(e.ID)
			Convey("Then we should receive an updated entity", func() {
				msg, err := n.Request("group.set", []byte(`{"id": `+id+`, "name":"foo"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "foo")
				So(err, ShouldEqual, nil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "foo")
			})
		})
	})

	Convey("Scenario: find datacenters", t, func() {
		setupTestSuite()
		Convey("Given groups exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of groups", func() {
				msg, _ := n.Request("group.find", []byte(`{}`), time.Second)
				list := []Entity{}
				json.Unmarshal(msg.Data, &list)
				So(len(list), ShouldEqual, 20)
			})
		})
	})
}
