package certificates

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agent/system"
	"github.com/r2dtools/agentintegration"
)

// Handler handles requests to the module
type Handler struct{}

// Handle handles request to the module
func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "issue":
		response, err = issue(request.Data)
	case "upload":
		response, err = upload(request.Data)
	case "storagecertnamelist":
		response, err = storageCertNameList(request.Data)
	case "storagecertdata":
		response, err = storageCertData(request.Data)
	case "storagecertupload":
		response, err = uploadCertToStorage(request.Data)
	case "storagecertremove":
		err = removeCertFromStorage(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func issue(data interface{}) (*agentintegration.Certificate, error) {
	var certData agentintegration.CertificateIssueRequestData
	err := mapstructure.Decode(data, &certData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	if err := system.GetPrivilege().IncreasePrivilege(); err != nil {
		logger.Error(fmt.Sprintf("certificate issue: increase privilege failed: %v", err))
	}

	defer system.GetPrivilege().DropPrivilege()
	certManager, err := GetCertificateManager()

	if err != nil {
		return nil, err
	}

	return certManager.Issue(certData)
}

func upload(data interface{}) (*agentintegration.Certificate, error) {
	var requestData agentintegration.CertificateUploadRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	if requestData.ServerName == "" {
		return nil, errors.New("domain name is missed")
	}

	if err := system.GetPrivilege().IncreasePrivilege(); err != nil {
		logger.Error(fmt.Sprintf("certificate issue: increase privilege failed: %v", err))
	}

	defer system.GetPrivilege().DropPrivilege()
	certManager, err := GetCertificateManager()

	if err != nil {
		return nil, err
	}

	return certManager.Upload(requestData.ServerName, requestData.WebServer, requestData.PemCertificate)
}

func storageCertNameList(data interface{}) (*agentintegration.StorageCertificateNameList, error) {
	certificateManager, err := GetCertificateManager()
	if err != nil {
		return nil, err
	}
	certList, err := certificateManager.GetStorageCertList()
	if err != nil {
		return nil, err
	}
	certNameList := agentintegration.StorageCertificateNameList{
		CertNameList: certList,
	}
	return &certNameList, nil
}

func storageCertData(data interface{}) (*agentintegration.Certificate, error) {
	certName, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid certificate name data is provided")
	}
	certificateManager, err := GetCertificateManager()
	if err != nil {
		return nil, err
	}

	return certificateManager.GetStorageCertData(certName)
}

func uploadCertToStorage(data interface{}) (*agentintegration.Certificate, error) {
	var requestData agentintegration.CertificateUploadRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	if requestData.CertName == "" {
		return nil, errors.New("certificate name is missed")
	}

	storage := GetDefaultCertStorage()
	_, _, err = storage.AddPemCertificate(requestData.CertName, requestData.PemCertificate)
	if err != nil {
		return nil, err
	}

	return storage.GetCertificate(requestData.CertName)
}

func removeCertFromStorage(data interface{}) error {
	certName, ok := data.(string)
	if !ok {
		return errors.New("invalid certificate name data is provided")
	}
	storage := GetDefaultCertStorage()

	return storage.RemoveCertificate(certName)
}
