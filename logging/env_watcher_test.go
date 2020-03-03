package logging_test

// Copyright (c) 2018, Arm Limited and affiliates.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


import (
	"fmt"

	. "github.com/armPelionEdge/wigwag-go-logger/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"
)

func setLogLevel(ll string) {
	Expect(ioutil.WriteFile("/tmp/log_level", []byte(ll), os.ModePerm)).Should(Succeed())
}

func setLogComponent(ll string) {
	Expect(ioutil.WriteFile("/tmp/log_component", []byte(ll), os.ModePerm)).Should(Succeed())
}

var _ = Describe("EnvWatcher", func() {
	BeforeEach(func() {
		os.Setenv(LogLevelEnvironmentVariable, "/tmp/log_level")
		os.Setenv(LogComponentEnvironmentVariable, "/tmp/log_component")

		go func() {
			for {
				<-time.After(time.Second)
				Log.Debugf("DEBUG")
				Log.Infof("INFO")
				Log.Noticef("NOTICE")
				Log.Warningf("WARNING")
				Log.Errorf("ERROR")
				Log.Criticalf("CRITICAL")
			}
		}()
	})

	Specify("It should change the log level depending on what the log level environment variable is set to", func() {
		// Initialize the log config in the applciation layer instead of get it from environment variable
		SetLoggingLevel("debug")

		SetLoggingComponent("WigWagGoLogger")

		// Set the component name
		setLogComponent("WigWagGoLogger")

		fmt.Println("Should see all levels")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		setLogLevel("asadf")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		setLogLevel("info")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		fmt.Println("Should see all levels except debug")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		setLogLevel("asfd")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		fmt.Println("Should see all levels except debug")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		setLogLevel("warning")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		fmt.Println("Should only see warning, error, critical")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))
	})

	Specify("It should change the log component depending on what the log component environment variable is set to", func() {
		// Initialize the log config in the applciation layer instead of get it from environment variable
		SetLoggingLevel("debug")

		SetLoggingComponent("WigWagGoLogger")

		fmt.Println("Should see component name with WigWagGoLogger")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))

		fmt.Println("Should see component name with WigWagGoLogger2")

		setLogComponent("WigWagGoLogger_update")

		<-time.After(time.Second * time.Duration(LogConfigSyncPeriodSeconds+1))
	})
})
