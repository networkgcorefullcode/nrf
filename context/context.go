// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/omec-project/nrf/factory"
	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/openapi/models"
)

var NrfNfProfile models.NfProfile
var Ipv4Address_aux string
var Transport_aux models.TransportProtocol
var Port_aux int32

func InitNrfContext() {
	config := factory.NrfConfig
	logger.InitLog.Infof("nrfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration

	NrfNfProfile.NfInstanceId = uuid.New().String()
	NrfNfProfile.NfType = models.NfType_NRF
	NrfNfProfile.NfStatus = models.NfStatus_REGISTERED

	serviceNameList := configuration.ServiceNameList
	NFServices := InitNFService(serviceNameList, config.Info.Version)
	NrfNfProfile.NfServices = &NFServices
}

func InitNFService(srvNameList []string, version string) []models.NfService {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	NFServices := make([]models.NfService, len(srvNameList))
	Ipv4Address_aux = factory.NrfConfig.GetSbiRegisterIP()
	Transport_aux = models.TransportProtocol_TCP
	Port_aux = int32(factory.NrfConfig.GetSbiPort())
	for index, nameString := range srvNameList {
		name := models.ServiceName(nameString)
		NFServices[index] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: []models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(factory.NrfConfig.GetSbiScheme()),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       &factory.SbiUri,
			IpEndPoints: []models.IpEndPoint{
				{
					Ipv4Address: &Ipv4Address_aux,
					Transport:   &Transport_aux,
					Port:        &Port_aux,
				},
			},
		}
	}
	return NFServices
}
