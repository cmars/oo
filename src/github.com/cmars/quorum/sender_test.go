/*
 * Copyright 2015 Casey Marshall
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package quorum_test

import (
	"sync"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/cmars/quorum"
)

type SenderSuite struct {
	sender quorum.MemSender
}

var _ = gc.Suite(&SenderSuite{})

func (s *SenderSuite) SetUpTest(c *gc.C) {
	s.sender = quorum.NewMemSender()
}

func (s *SenderSuite) TestValidate(c *gc.C) {
	err := s.sender.ValidateRecipient("alice")
	c.Assert(err, gc.ErrorMatches, "not found")
	s.sender.Register("alice", func(b quorum.Ballot) error { panic("not called") })
	err = s.sender.ValidateRecipient("alice")
	c.Assert(err, gc.IsNil)
}

func (s *SenderSuite) TestSend(c *gc.C) {
	var err error
	var mu sync.Mutex
	var recipients []string

	f := func(b quorum.Ballot) error {
		mu.Lock()
		c.Logf("ballot for %q", b.Recipient)
		recipients = append(recipients, b.Recipient)
		mu.Unlock()
		return nil
	}
	s.sender.Register("alice", f)
	s.sender.Register("bob", f)
	err = s.sender.Send(quorum.Ballot{
		Recipient: "alice",
	})
	c.Assert(err, jc.ErrorIsNil)
	err = s.sender.Send(quorum.Ballot{
		Recipient: "bob",
	})
	c.Assert(err, jc.ErrorIsNil)
	err = s.sender.Send(quorum.Ballot{
		Recipient: "dave",
	})
	c.Assert(err, gc.ErrorMatches, "not found")

	s.sender.Close()

	mu.Lock()
	defer mu.Unlock()
	c.Assert(recipients, jc.SameContents, []string{"alice", "bob"})
}
