package logging_test

import (
	"fmt"

	. "wigwag-go-logger/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"
)

func setLogLevel(ll string) {
	Expect(ioutil.WriteFile("/tmp/log_level", []byte(ll), os.ModePerm)).Should(Succeed())
}

var _ = Describe("EnvWatcher", func() {
	BeforeEach(func() {
		os.Setenv(LogLevelEnvironmentVariable, "/tmp/log_level")

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
		// Set the component name
		SetLoggingComponent("WigWagGoLogger")
		// Set it to the highest level first so we can watch it adjust
		SetLoggingLevel("debug")

		fmt.Println("Should see all levels")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		setLogLevel("asadf")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		setLogLevel("info")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		fmt.Println("Should see all levels except debug")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		setLogLevel("asfd")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		fmt.Println("Should see all levels except debug")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		setLogLevel("warning")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))

		fmt.Println("Should only see warning, error, critical")

		<-time.After(time.Second * time.Duration(LogLevelSyncPeriodSeconds+1))
	})
})
