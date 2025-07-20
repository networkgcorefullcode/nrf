// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package management

import (
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/omec-project/nrf/dbadapter"
	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/nrf/util"
	"github.com/omec-project/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNrfInfo() *models.NrfInfo {
	// init
	var nrfinfo models.NrfInfo

	nrfinfo.ServedUdrInfo = getUdrInfo()
	nrfinfo.ServedUdmInfo = getUdmInfo()
	nrfinfo.ServedAusfInfo = getAusfInfo()
	nrfinfo.ServedAmfInfo = getAmfInfo()
	nrfinfo.ServedSmfInfo = getSmfInfo()
	nrfinfo.ServedUpfInfo = getUpfInfo()
	nrfinfo.ServedPcfInfo = getPcfInfo()
	nrfinfo.ServedBsfInfo = getBsfInfo()
	nrfinfo.ServedChfInfo = getChfInfo()

	return &nrfinfo
}

func getUdrInfo() *map[string]models.NrfInfoServedUdrInfoValue {
	servedUdrInfo := make(map[string]models.NrfInfoServedUdrInfoValue)
	var UDRProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "UDR"}

	UDR, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	UDRStruct, err := util.Decode(UDR, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(UDRStruct); i++ {
		err := mapstructure.Decode(UDRStruct[i], &UDRProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedUdrInfo[index] = models.NrfInfoServedUdrInfoValue{
			UdrInfo: UDRProfile.UdrInfo,
		}
	}
	return &servedUdrInfo
}

func getUdmInfo() *map[string]models.NrfInfoServedUdmInfoValue {
	servedUdmInfo := make(map[string]models.NrfInfoServedUdmInfoValue)
	var UDMProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "UDM"}

	UDM, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	UDMStruct, err := util.Decode(UDM, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(UDMStruct); i++ {
		err := mapstructure.Decode(UDMStruct[i], &UDMProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedUdmInfo[index] = models.NrfInfoServedUdmInfoValue{
			UdmInfo: UDMProfile.UdmInfo,
		}
	}
	return &servedUdmInfo
}

func getAusfInfo() *map[string]models.NrfInfoServedAusfInfoValue {
	servedAusfInfo := make(map[string]models.NrfInfoServedAusfInfoValue)
	var AUSFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "AUSF"}

	AUSF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	AUSFStruct, err := util.Decode(AUSF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(AUSFStruct); i++ {
		err := mapstructure.Decode(AUSFStruct[i], &AUSFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedAusfInfo[index] = models.NrfInfoServedAusfInfoValue{
			AusfInfo: AUSFProfile.AusfInfo,
		}
	}
	return &servedAusfInfo
}

func getAmfInfo() *map[string]models.NrfInfoServedAmfInfoValue {
	servedAmfInfo := make(map[string]models.NrfInfoServedAmfInfoValue)
	var AMFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "AMF"}

	AMF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	AMFStruct, err := util.Decode(AMF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(AMFStruct); i++ {
		err := mapstructure.Decode(AMFStruct[i], &AMFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedAmfInfo[index] = models.NrfInfoServedAmfInfoValue{
			AmfInfo: AMFProfile.AmfInfo,
		}
	}
	return &servedAmfInfo
}

func getSmfInfo() *map[string]models.NrfInfoServedSmfInfoValue {
	servedSmfInfo := make(map[string]models.NrfInfoServedSmfInfoValue)
	var SMFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "SMF"}

	SMF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	SMFStruct, err := util.Decode(SMF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(SMFStruct); i++ {
		err := mapstructure.Decode(SMFStruct[i], &SMFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedSmfInfo[index] = models.NrfInfoServedSmfInfoValue{
			SmfInfo: SMFProfile.SmfInfo,
		}
	}
	return &servedSmfInfo
}

func getUpfInfo() *map[string]models.NrfInfoServedUpfInfoValue {
	servedUpfInfo := make(map[string]models.NrfInfoServedUpfInfoValue)
	var UPFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "UPF"}

	UPF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	UPFStruct, err := util.Decode(UPF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(UPFStruct); i++ {
		err := mapstructure.Decode(UPFStruct[i], &UPFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedUpfInfo[index] = models.NrfInfoServedUpfInfoValue{
			UpfInfo: UPFProfile.UpfInfo,
		}
	}
	return &servedUpfInfo
}

func getPcfInfo() *map[string]models.NrfInfoServedPcfInfoValue {
	servedPcfInfo := make(map[string]models.NrfInfoServedPcfInfoValue)
	var PCFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "PCF"}

	PCF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	PCFStruct, err := util.Decode(PCF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(PCFStruct); i++ {
		err := mapstructure.Decode(PCFStruct[i], &PCFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedPcfInfo[index] = models.NrfInfoServedPcfInfoValue{
			PcfInfo: PCFProfile.PcfInfo,
		}
	}
	return &servedPcfInfo
}

func getBsfInfo() *map[string]models.NrfInfoServedBsfInfoValue {
	servedBsfInfo := make(map[string]models.NrfInfoServedBsfInfoValue)
	var BSFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "BSF"}

	BSF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	BSFStruct, err := util.Decode(BSF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(BSFStruct); i++ {
		err := mapstructure.Decode(BSFStruct[i], &BSFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedBsfInfo[index] = models.NrfInfoServedBsfInfoValue{
			BsfInfo: BSFProfile.BsfInfo,
		}
	}
	return &servedBsfInfo
}

func getChfInfo() *map[string]models.NrfInfoServedChfInfoValue {
	servedChfInfo := make(map[string]models.NrfInfoServedChfInfoValue)
	var CHFProfile models.NfProfile

	collName := "NfProfile"
	filter := bson.M{"nfType": "CHF"}

	CHF, _ := dbadapter.DBClient.RestfulAPIGetMany(collName, filter)
	CHFStruct, err := util.Decode(CHF, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for i := 0; i < len(CHFStruct); i++ {
		err := mapstructure.Decode(CHFStruct[i], &CHFProfile)
		if err != nil {
			panic(err)
		}
		index := strconv.Itoa(i)
		servedChfInfo[index] = models.NrfInfoServedChfInfoValue{
			ChfInfo: CHFProfile.ChfInfo,
		}
	}
	return &servedChfInfo
}

// DecodeNfProfile - Only support []map[string]interface to []models.NfProfile
func DecodeNfProfile(source interface{}, format string) (models.NfProfile, error) {
	var target models.NfProfile

	// config mapstruct
	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(format, data.(string))
		}
		return data, nil
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: stringToDateTimeHook,
		Result:     &target,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return target, err
	}

	// Decode result to NfProfile structure
	err = decoder.Decode(source)
	if err != nil {
		return target, err
	}
	return target, nil
}
