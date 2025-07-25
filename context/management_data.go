// SPDX-FileCopyrightText: 2025 Intel Corporation
// SPDX-FileCopyrightText: 2025 Canonical Ltd
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/omec-project/nrf/dbadapter"
	"github.com/omec-project/nrf/factory"
	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/nrf/polling"
	"github.com/omec-project/openapi"
	"github.com/omec-project/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

const NRF_NFINST_RES_URI_PREFIX = factory.NRF_NFM_RES_URI_PREFIX + "/nf-instances/"

// Generates a random int between 0 and 99
func GenerateRandomNumber() (int, error) {
	maximum := big.NewInt(100)
	randomNumber, err := rand.Int(rand.Reader, maximum)
	if err != nil {
		return 0, err
	}
	return int(randomNumber.Int64()), nil
}

func NnrfNFManagementDataModel(nf *models.NfProfile, nfprofile models.NfProfile) error {
	if nfprofile.NfInstanceId == "" {
		return fmt.Errorf("NfInstanceId field is required")
	}
	nf.NfInstanceId = nfprofile.NfInstanceId

	if nfprofile.NfType == "" {
		return fmt.Errorf("NfType field is required")
	}
	nf.NfType = nfprofile.NfType

	if nfprofile.NfStatus == "" {
		return fmt.Errorf("NfStatus field is required")
	}
	nf.NfStatus = nfprofile.NfStatus

	nfPlmnList, err := buildNfProfilePlmnList(nfprofile.PlmnList)
	if err != nil {
		return err
	}

	nnrfNFManagementCondition(nf, nfprofile)
	nf.PlmnList = &nfPlmnList
	nnrfNFManagementOption(nf, nfprofile)

	return nil
}

func buildNfProfilePlmnList(nfProvidedPlmnList *[]models.PlmnId) ([]models.PlmnId, error) {
	// NF provided a list of supported PLMNs
	if nfProvidedPlmnList != nil && len(*nfProvidedPlmnList) != 0 {
		return *nfProvidedPlmnList, nil
	}
	// NF did not provide supported PLMNs: fetch from webconsole
	logger.ManagementLog.Warnln("PLMN config not provided by NF, using supported PLMNs from webconsole")
	supportedPlmnList, err := polling.FetchPlmnConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PLMN config from webconsole: %v", err)
	}
	logger.ManagementLog.Debugf("Fetched PLMN list from webconsole: %+v", supportedPlmnList)
	if len(supportedPlmnList) == 0 {
		return nil, fmt.Errorf("PLMN config not provided by NF and no local PLMN config available")
	}
	return supportedPlmnList, nil
}

func SetsubscriptionId() string {
	x, err := GenerateRandomNumber()
	if err != nil {
		logger.ManagementLog.Error(err)
	}
	return strconv.Itoa(x)
}

func nnrfNFManagementCondition(nf *models.NfProfile, nfprofile models.NfProfile) {
	// HeartBeatTimer
	if !factory.NrfConfig.Configuration.NfProfileExpiryEnable {
		// setting 1day keepAliveTimer value
		factory.NrfConfig.Configuration.NfKeepAliveTime = 24 * 60 * 60
	} else if factory.NrfConfig.Configuration.NfKeepAliveTime == 0 {
		logger.ManagementLog.Infoln("NfProfileExpiryEnable: true but keepAliveTime: 0, setting default keepAliveTimer: 60 sec")
		factory.NrfConfig.Configuration.NfKeepAliveTime = 60
	}
	nf.HeartBeatTimer = factory.NrfConfig.Configuration.NfKeepAliveTime
	logger.ManagementLog.Infof("HeartBeat Timer value: %v sec", nf.HeartBeatTimer)

	// fqdn
	if nfprofile.Fqdn != "" {
		nf.Fqdn = nfprofile.Fqdn
	}
	// interPlmnFqdn
	if nfprofile.InterPlmnFqdn != "" {
		nf.InterPlmnFqdn = nfprofile.InterPlmnFqdn
	}
	// ipv4Addresses
	if nfprofile.Ipv4Addresses != nil {
		// fmt.Println("NsiList")
		a := make([]string, len(nfprofile.Ipv4Addresses))
		copy(a, nfprofile.Ipv4Addresses)
		nf.Ipv4Addresses = a
	}
	// ipv6Addresses
	if nfprofile.Ipv6Addresses != nil {
		// fmt.Println("NsiList")
		a := make([]string, len(nfprofile.Ipv6Addresses))
		copy(a, nfprofile.Ipv6Addresses)
		nf.Ipv6Addresses = a
	}
}

func nnrfNFManagementOption(nf *models.NfProfile, nfprofile models.NfProfile) {
	// sNssais
	if nfprofile.SNssais != nil {
		// fmt.Println("SNssais")
		a := make([]models.Snssai, len(*nfprofile.SNssais))
		copy(a, *nfprofile.SNssais)
		nf.SNssais = &a
	}

	// nsiList
	if nfprofile.NsiList != nil {
		// fmt.Println("NsiList")
		a := make([]string, len(nfprofile.NsiList))
		copy(a, nfprofile.NsiList)
		nf.NsiList = a
	}
	// allowedPlmns
	if nfprofile.AllowedPlmns != nil {
		a := make([]models.PlmnId, len(*nfprofile.AllowedPlmns))
		copy(a, *nfprofile.AllowedPlmns)
		nf.AllowedPlmns = &a
	}
	// allowedNfTypes
	if nfprofile.AllowedNfTypes != nil {
		a := make([]models.NfType, len(nfprofile.AllowedNfTypes))
		copy(a, nfprofile.AllowedNfTypes)
		nf.AllowedNfTypes = a
	}
	// allowedNfDomains
	if nfprofile.AllowedNfDomains != nil {
		a := make([]string, len(nfprofile.AllowedNfDomains))
		copy(a, nfprofile.AllowedNfDomains)
		nf.AllowedNfDomains = a
	}

	// allowedNssais
	if nfprofile.AllowedNssais != nil {
		// fmt.Println("SNssais")
		a := make([]models.Snssai, len(*nfprofile.AllowedNssais))
		copy(a, *nfprofile.AllowedNssais)
		nf.AllowedNssais = &a
	}
	// Priority
	if nfprofile.Priority > 0 && nfprofile.Priority <= 65535 {
		nf.Priority = nfprofile.Priority
	}
	// Capacity
	if nfprofile.Capacity > 0 && nfprofile.Capacity <= 65535 {
		nf.Capacity = nfprofile.Capacity
	}
	// Load
	if nfprofile.Load > 0 && nfprofile.Load <= 100 {
		nf.Load = nfprofile.Load
	}
	// Locality
	if nfprofile.Locality != "" {
		nf.Locality = nfprofile.Locality
	}
	// udrInfo
	if nfprofile.UdrInfo != nil {
		var a models.UdrInfo

		if nfprofile.UdrInfo.GroupId != "" {
			a.GroupId = nfprofile.UdrInfo.GroupId
		}

		if nfprofile.UdrInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.UdrInfo.SupiRanges
		}

		if nfprofile.UdrInfo.GpsiRanges != nil {
			a.GpsiRanges = nfprofile.UdrInfo.GpsiRanges
		}

		if nfprofile.UdrInfo.ExternalGroupIdentifiersRanges != nil {
			a.ExternalGroupIdentifiersRanges = nfprofile.UdrInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdrInfo.SupportedDataSets != nil {
			a.SupportedDataSets = nfprofile.UdrInfo.SupportedDataSets
		}

		nf.UdrInfo = &a
	}
	// udmInfo
	if nfprofile.UdmInfo != nil {
		var a models.UdmInfo

		if nfprofile.UdmInfo.GroupId != "" {
			a.GroupId = nfprofile.UdmInfo.GroupId
		}

		if nfprofile.UdmInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.UdmInfo.SupiRanges
		}

		if nfprofile.UdmInfo.GpsiRanges != nil {
			a.GpsiRanges = nfprofile.UdmInfo.GpsiRanges
		}

		if nfprofile.UdmInfo.ExternalGroupIdentifiersRanges != nil {
			a.ExternalGroupIdentifiersRanges = nfprofile.UdmInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdmInfo.RoutingIndicators != nil {
			a.RoutingIndicators = nfprofile.UdmInfo.RoutingIndicators
		}

		nf.UdmInfo = &a
	}
	// ausfInfo
	if nfprofile.AusfInfo != nil {
		var a models.AusfInfo

		if nfprofile.AusfInfo.GroupId != "" {
			a.GroupId = nfprofile.AusfInfo.GroupId
		}

		if nfprofile.AusfInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.AusfInfo.SupiRanges
		}

		if nfprofile.AusfInfo.RoutingIndicators != nil {
			a.RoutingIndicators = nfprofile.AusfInfo.RoutingIndicators
		}

		nf.AusfInfo = &a
	}
	// amfInfo
	if nfprofile.AmfInfo != nil {
		var a models.AmfInfo

		if nfprofile.AmfInfo.AmfSetId != "" {
			a.AmfSetId = nfprofile.AmfInfo.AmfSetId
		}

		if nfprofile.AmfInfo.AmfRegionId != "" {
			a.AmfRegionId = nfprofile.AmfInfo.AmfRegionId
		}

		if nfprofile.AmfInfo.GuamiList != nil {
			a.GuamiList = nfprofile.AmfInfo.GuamiList
		}

		if nfprofile.AmfInfo.TaiList != nil {
			a.TaiList = nfprofile.AmfInfo.TaiList
		}

		if nfprofile.AmfInfo.TaiRangeList != nil {
			a.TaiRangeList = nfprofile.AmfInfo.TaiRangeList
		}

		if nfprofile.AmfInfo.BackupInfoAmfFailure != nil {
			a.BackupInfoAmfFailure = nfprofile.AmfInfo.BackupInfoAmfFailure
		}

		if nfprofile.AmfInfo.BackupInfoAmfRemoval != nil {
			a.BackupInfoAmfRemoval = nfprofile.AmfInfo.BackupInfoAmfRemoval
		}

		if nfprofile.AmfInfo.N2InterfaceAmfInfo != nil {
			a.N2InterfaceAmfInfo = nfprofile.AmfInfo.N2InterfaceAmfInfo
		}
		nf.AmfInfo = &a
	}
	// smfInfo
	if nfprofile.SmfInfo != nil {
		var a models.SmfInfo

		if nfprofile.SmfInfo.SNssaiSmfInfoList != nil {
			logger.ManagementLog.Debugln("Setting SNssaiSmfInfoList in SmfInfo")
			a.SNssaiSmfInfoList = nfprofile.SmfInfo.SNssaiSmfInfoList
		}
		if nfprofile.SmfInfo.TaiList != nil {
			logger.ManagementLog.Debugln("Setting TaiList in SmfInfo")
			a.TaiList = nfprofile.SmfInfo.TaiList
		}
		if nfprofile.SmfInfo.TaiRangeList != nil {
			logger.ManagementLog.Debugln("Setting TaiRangeList in SmfInfo")
			a.TaiRangeList = nfprofile.SmfInfo.TaiRangeList
		}
		if nfprofile.SmfInfo.PgwFqdn != "" {
			logger.ManagementLog.Debugf("Setting PgwFqdn in SmfInfo: %s", nfprofile.SmfInfo.PgwFqdn)
			a.PgwFqdn = nfprofile.SmfInfo.PgwFqdn
		}
		if nfprofile.SmfInfo.AccessType != nil {
			logger.ManagementLog.Debugln("Setting AccessType in SmfInfo")
			a.AccessType = nfprofile.SmfInfo.AccessType
		}
		nf.SmfInfo = &a
	}
	// upfInfo
	if nfprofile.UpfInfo != nil {
		var a models.UpfInfo

		if nfprofile.UpfInfo.SNssaiUpfInfoList != nil {
			a.SNssaiUpfInfoList = nfprofile.UpfInfo.SNssaiUpfInfoList
		}
		if nfprofile.UpfInfo.SmfServingArea != nil {
			a.SmfServingArea = nfprofile.UpfInfo.SmfServingArea
		}
		if nfprofile.UpfInfo.InterfaceUpfInfoList != nil {
			a.InterfaceUpfInfoList = nfprofile.UpfInfo.InterfaceUpfInfoList
		}

		a.IwkEpsInd = nfprofile.UpfInfo.IwkEpsInd

		nf.UpfInfo = &a
	}
	// pcfInfo
	if nfprofile.PcfInfo != nil {
		var a models.PcfInfo

		if nfprofile.PcfInfo.DnnList != nil {
			a.DnnList = nfprofile.PcfInfo.DnnList
		}
		if nfprofile.PcfInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.PcfInfo.SupiRanges
		}
		if nfprofile.PcfInfo.RxDiamHost != "" {
			a.RxDiamHost = nfprofile.PcfInfo.RxDiamHost
		}
		if nfprofile.PcfInfo.RxDiamRealm != "" {
			a.RxDiamRealm = nfprofile.PcfInfo.RxDiamRealm
		}
		nf.PcfInfo = &a
	}
	// bsfInfo
	if nfprofile.BsfInfo != nil {
		var a models.BsfInfo

		if nfprofile.BsfInfo.DnnList != nil {
			a.DnnList = nfprofile.BsfInfo.DnnList
		}
		if nfprofile.BsfInfo.IpDomainList != nil {
			a.IpDomainList = nfprofile.BsfInfo.IpDomainList
		}
		if nfprofile.BsfInfo.Ipv4AddressRanges != nil {
			b := make([]models.Ipv4AddressRange, len(*nfprofile.BsfInfo.Ipv4AddressRanges))
			for i := 0; i < len(*nfprofile.BsfInfo.Ipv4AddressRanges); i++ {
				b[i].Start = strconv.Itoa(int(Ipv4ToInt((*nfprofile.BsfInfo.Ipv4AddressRanges)[i].Start)))
				b[i].End = strconv.Itoa(int(Ipv4ToInt((*nfprofile.BsfInfo.Ipv4AddressRanges)[i].End)))
			}
			a.Ipv4AddressRanges = &b
		}
		if nfprofile.BsfInfo.Ipv6PrefixRanges != nil {
			b := make([]models.Ipv6PrefixRange, len(*nfprofile.BsfInfo.Ipv6PrefixRanges))
			for i := 0; i < len(*nfprofile.BsfInfo.Ipv6PrefixRanges); i++ {
				b[i].Start = Ipv6ToInt(((*nfprofile.BsfInfo.Ipv6PrefixRanges)[i].Start)).String()
				b[i].End = Ipv6ToInt(((*nfprofile.BsfInfo.Ipv6PrefixRanges)[i].End)).String()
			}
			a.Ipv6PrefixRanges = &b
		}
		nf.BsfInfo = &a
	}
	// chfInfo
	if nfprofile.ChfInfo != nil {
		var a models.ChfInfo

		if nfprofile.ChfInfo.SupiRangeList != nil {
			a.SupiRangeList = nfprofile.ChfInfo.SupiRangeList
		}
		if nfprofile.ChfInfo.GpsiRangeList != nil {
			a.GpsiRangeList = nfprofile.ChfInfo.GpsiRangeList
		}
		if nfprofile.ChfInfo.PlmnRangeList != nil {
			a.PlmnRangeList = nfprofile.ChfInfo.PlmnRangeList
		}
		nf.ChfInfo = &a
	}
	// nrfInfo
	if nfprofile.NrfInfo != nil {
		nf.NrfInfo = nfprofile.NrfInfo
	}
	// recoveryTime
	if nfprofile.RecoveryTime != nil {
		// Update when restart (Setting by NF itself)
		nf.RecoveryTime = nfprofile.RecoveryTime
	}

	// nfServicePersistence
	if nfprofile.NfServicePersistence {
		nf.NfServicePersistence = true
	} else {
		nf.NfServicePersistence = false
	}

	// nfServices
	if nfprofile.NfServices != nil {
		a := make([]models.NfService, len(*nfprofile.NfServices))
		copy(a, *nfprofile.NfServices)
		nf.NfServices = &a
	}
	// fill the NfServiceList if NfServices is set
	if nfprofile.NfServices != nil && nfprofile.NfServiceList == nil {
		logger.ManagementLog.Debugln("NfServiceList is nil, setting NfServiceList from NfServices")

		map_service := make(map[string]models.NfService, len(*nfprofile.NfServices))
		for _, nfService := range *nfprofile.NfServices {
			map_service[string(nfService.ServiceInstanceId)] = nfService
		}
		nf.NfServiceList = &map_service
	}

	// fill the NfServices if NfServiceList is set
	if nfprofile.NfServiceList != nil && nfprofile.NfServices == nil {
		logger.ManagementLog.Debugln("NfServices is nil, setting NfServices from NfServiceList")

		var nfServices []models.NfService
		for _, nfService := range *nfprofile.NfServiceList {
			nfServices = append(nfServices, nfService)
		}
		nf.NfServices = &nfServices
	}
	logger.ManagementLog.Debugln("finish the function nnrfNFManagementOption")
}

func GetNfInstanceURI(nfInstID string) string {
	return factory.NrfConfig.GetSbiUri() + NRF_NFINST_RES_URI_PREFIX + nfInstID
}

func SetLocationHeader(nfprofile models.NfProfile) string {
	var modifyUL UriList
	var locationHeader []string

	// set nfprofile location
	locationHeader = append(locationHeader, GetNfInstanceURI(nfprofile.NfInstanceId))

	collName := "urilist"
	nfType := nfprofile.NfType
	filter := bson.M{"nfType": nfType}

	ul, _ := dbadapter.DBClient.RestfulAPIGetOne(collName, filter)

	var originalUL UriList
	err := mapstructure.Decode(ul, &originalUL)
	if err != nil {
		panic(err)
	}

	// obtain location header = NF URI
	nnrfUriList(&originalUL, &modifyUL, locationHeader)
	modifyUL.NfType = nfprofile.NfType

	tmp, err := json.Marshal(modifyUL)
	if err != nil {
		logger.ManagementLog.Error(err)
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	if ok, _ := dbadapter.DBClient.RestfulAPIPutOne(collName, filter, putData); ok {
		logger.ManagementLog.Info("urilist update")
	} else {
		logger.ManagementLog.Info("urilist create")
	}

	return locationHeader[0]
}

func setUriListByFilter(filter bson.M, uriList *[]string) {
	filterNfTypeResultsRaw, _ := dbadapter.DBClient.RestfulAPIGetMany("Subscriptions", filter)
	var filterNfTypeResults []models.NrfSubscriptionData
	err := openapi.Convert(filterNfTypeResultsRaw, &filterNfTypeResults)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for _, subscr := range filterNfTypeResults {
		*uriList = append(*uriList, subscr.NfStatusNotificationUri)
	}
}

func nnrfUriList(originalUL *UriList, UL *UriList, location []string) {
	var b *Links
	var flag bool
	var c []models.Link
	flag = true
	b = new(Links)
	size := len(location) + len(originalUL.Link.Item)

	// check duplicate
	for _, item := range originalUL.Link.Item {
		if item.Href == location[0] {
			flag = false
			break
		}
	}

	if flag {
		c = make([]models.Link, size)
		copy(c, originalUL.Link.Item)
		for i, loc := range location {
			c[len(originalUL.Link.Item)+i].Href = loc
		}
	} else {
		c = make([]models.Link, size-1)
		copy(c, originalUL.Link.Item)
	}

	b.Item = c
	UL.Link = *b
}

func GetNotificationUri(nfProfile models.NfProfile) []string {
	var uriList []string

	// nfTypeCond
	nfTypeCond := bson.M{
		"subscrCond": bson.M{
			"nfType": nfProfile.NfType,
		},
	}
	setUriListByFilter(nfTypeCond, &uriList)

	// NfInstanceIdCond
	nfInstanceIDCond := bson.M{
		"subscrCond": bson.M{
			"nfInstanceId": nfProfile.NfInstanceId,
		},
	}
	setUriListByFilter(nfInstanceIDCond, &uriList)

	// ServiceNameCond
	if nfProfile.NfServices != nil {
		var ServiceNameCond bson.M
		var serviceNames bson.A
		for _, nfService := range *nfProfile.NfServices {
			serviceNames = append(serviceNames, string(nfService.ServiceName))
		}
		ServiceNameCond = bson.M{
			"subscrCond.serviceName": bson.M{
				"$in": serviceNames,
			},
		}
		setUriListByFilter(ServiceNameCond, &uriList)
	}

	// AmfCond
	if nfProfile.AmfInfo != nil {
		amfCond := bson.M{
			"subscrCond": bson.M{
				"amfSetId":    (*nfProfile.AmfInfo).AmfSetId,
				"amfRegionId": (*nfProfile.AmfInfo).AmfRegionId,
			},
		}
		setUriListByFilter(amfCond, &uriList)
	}

	// GuamiListCond
	if nfProfile.AmfInfo != nil {
		var guamiListFilter bson.M
		if (*nfProfile.AmfInfo).GuamiList != nil {
			var guamiListBsonArray bson.A
			for _, guami := range *(*nfProfile.AmfInfo).GuamiList {
				tmp, err := json.Marshal(guami)
				if err != nil {
					logger.ManagementLog.Error(err)
				}
				guamiMarshal := bson.M{}
				err = json.Unmarshal(tmp, &guamiMarshal)
				if err != nil {
					logger.ManagementLog.Error(err)
				}

				guamiListBsonArray = append(guamiListBsonArray, bson.M{"subscrCond": bson.M{"$elemMatch": guamiMarshal}})
			}
			guamiListFilter = bson.M{
				"$or": guamiListBsonArray,
			}
		}
		setUriListByFilter(guamiListFilter, &uriList)
	}

	// NetworkSliceCond
	if nfProfile.SNssais != nil {
		var networkSliceFilter bson.M
		var snssaisBsonArray bson.A
		for _, snssai := range *nfProfile.SNssais {
			tmp, err := json.Marshal(snssai)
			if err != nil {
				logger.ManagementLog.Error(err)
			}
			snssaiMarshal := bson.M{}
			err = json.Unmarshal(tmp, &snssaiMarshal)
			if err != nil {
				logger.ManagementLog.Error(err)
			}

			snssaisBsonArray = append(snssaisBsonArray, bson.M{"subscrCond": bson.M{"$elemMatch": snssaiMarshal}})
		}

		var nsiListBsonArray bson.A
		if nfProfile.NsiList != nil {
			for _, nsi := range nfProfile.NsiList {
				nsiListBsonArray = append(nsiListBsonArray, nsi)
			}
		}

		if nsiListBsonArray != nil {
			networkSliceFilter = bson.M{
				"$and": bson.A{
					bson.M{
						"subscrCond.nsiList": bson.M{
							"$in": nsiListBsonArray,
						},
					},
					bson.M{
						"$or": snssaisBsonArray,
					},
				},
			}
		} else {
			networkSliceFilter = bson.M{
				"$and": bson.A{
					bson.M{
						"$or": snssaisBsonArray,
					},
				},
			}
		}
		setUriListByFilter(networkSliceFilter, &uriList)
	}

	// NfGroupCond
	if nfProfile.UdrInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.UdrInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.UdmInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.UdmInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.AusfInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.AusfInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	}

	return uriList
}

func NnrfUriListLimit(originalUL *UriList, limit int) {
	// response limit

	if limit < len(originalUL.Link.Item) {
		b := new(Links)
		c := make([]models.Link, limit)
		copy(c, originalUL.Link.Item[:limit])
		b.Item = c
		originalUL.Link = *b
	}
}
