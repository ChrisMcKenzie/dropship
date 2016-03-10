// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package dropship

// Locker is an interface that allows you to block execution of
// another process across a set of machines.
type Locker interface {
	// Acquire takes a shutdownCh and return a lock chan and error
	//
	// the lock chan can be used to block the process until the lock
	// has been acquired and the chan receives.
	Acquire(<-chan struct{}) (<-chan struct{}, error)
	// Release will release the lock allowing for other processes to
	// acquire.
	Release() error
}
