/*
Copyright SecureKey Technologies Inc. All Rights Reserved.


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

package integration

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/util"
)

const (
	chainCodeName = "install"
	chainCodePath = "github.com/example_cc"
)

func TestChaincodeInstal(t *testing.T) {

	testSetup := &BaseSetupImpl{
		ConfigFile:      "../fixtures/config/config_test.yaml",
		ChainID:         "mychannel",
		ChannelConfig:   "../fixtures/channel/mychannel.tx",
		ConnectEventHub: true,
	}

	if err := testSetup.Initialize(); err != nil {
		t.Fatalf(err.Error())
	}

	testChaincodeInstallUsingChaincodePath(t, testSetup)

	testChaincodeInstallUsingChaincodePackage(t, testSetup)
}

// Test chaincode install using chaincodePath to create chaincodePackage
func testChaincodeInstallUsingChaincodePath(t *testing.T, testSetup *BaseSetupImpl) {
	chainCodeVersion := getRandomCCVersion()

	// Install and Instantiate Events CC
	// Retrieve installed chaincodes
	client := testSetup.Client

	if err := testSetup.InstallCC(chainCodeName, chainCodePath, chainCodeVersion, nil); err != nil {
		t.Fatalf("installCC return error: %v", err)
	}
	chaincodeQueryResponse, err := client.QueryInstalledChaincodes(testSetup.Chain.GetPrimaryPeer())
	if err != nil {
		t.Fatalf("QueryInstalledChaincodes return error: %v", err)
	}
	ccFound := false
	for _, chaincode := range chaincodeQueryResponse.Chaincodes {
		if chaincode.Name == chainCodeName && chaincode.Path == chainCodePath && chaincode.Version == chainCodeVersion {
			fmt.Printf("Found chaincode: %s\n", chaincode)
			ccFound = true
		}
	}

	if !ccFound {
		t.Fatalf("Failed to retrieve installed chaincode.")
	}
	//Install same chaincode again, should fail
	err = testSetup.InstallCC(chainCodeName, chainCodePath, chainCodeVersion, nil)
	if err == nil {
		t.Fatalf("install same chaincode didn't return error")
	}
	if strings.Contains(err.Error(), "chaincodes/install.v"+chainCodeVersion+" exists") {
		t.Fatalf("install same chaincode didn't return the correct error")
	}
}

// Test chaincode install using chaincodePackage[byte]
func testChaincodeInstallUsingChaincodePackage(t *testing.T, testSetup *BaseSetupImpl) {

	chainCodeVersion := getRandomCCVersion()
	util.ChangeGOPATHToDeploy(testSetup.GetDeployPath())
	chaincodePackage, err := fabricClient.PackageCC(chainCodePath, "")
	util.ResetGOPATH()
	if err != nil {
		t.Fatalf("PackageCC return error: %s", err)
	}

	err = testSetup.InstallCC("install", "github.com/example_cc_pkg", chainCodeVersion, chaincodePackage)
	if err != nil {
		t.Fatalf("installCC return error: %v", err)
	}
	//Install same chaincode again, should fail
	err = testSetup.InstallCC("install", chainCodePath, chainCodeVersion, chaincodePackage)
	if err == nil {
		t.Fatalf("install same chaincode didn't return error")
	}
	if strings.Contains(err.Error(), "chaincodes/install.v"+chainCodeVersion+" exists") {
		t.Fatalf("install same chaincode didn't return the correct error")
	}
}

func getRandomCCVersion() string {
	rand.Seed(time.Now().UnixNano())
	return "v0" + strconv.Itoa(rand.Intn(10000000))
}