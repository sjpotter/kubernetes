/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"os"
	"sync"
	"time"
)

// FakeOS mocks out certain OS calls to avoid perturbing the filesystem
// on the test machine.
// TODO(yifan): Use go mock to replace this.
type FakeOS struct {
	HostName string
	Removes  []string
	sync.Mutex
}

func NewFakeOS() *FakeOS {
	return &FakeOS{Removes: []string{}}
}

// Mkdir is a fake call that just returns nil.
func (FakeOS) Mkdir(path string, perm os.FileMode) error {
	return nil
}

// Symlink is a fake call that just returns nil.
func (FakeOS) Symlink(oldname string, newname string) error {
	return nil
}

// Remove is a fake call that returns nil.
func (f *FakeOS) Remove(path string) error {
	f.Lock()
	defer f.Unlock()
	f.Removes = append(f.Removes, path)
	return nil
}

// Create is a fake call that returns nil.
func (FakeOS) Create(path string) (*os.File, error) {
	return nil, nil
}

// Hostname is a fake call that returns nil.
func (f *FakeOS) Hostname() (name string, err error) {
	return f.HostName, nil
}

// Chtimes is a fake call that returns nil.
func (FakeOS) Chtimes(path string, atime time.Time, mtime time.Time) error {
	return nil
}

// Pipe is a fake call that returns nil.
func (FakeOS) Pipe() (r *os.File, w *os.File, err error) {
	return nil, nil, nil
}
