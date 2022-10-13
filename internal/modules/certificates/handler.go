package certificates

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agent/pkg/router"
	"github.com/r2dtools/agentintegration"
)

type Handler struct {
	certificateManager *CertificateManager
	logger             logger.LoggerInterface
}

func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "issue":
		response, err = h.issueCertificateToDomain(request.Data)
	case "upload":
		response, err = h.uploadCertificateToDomain(request.Data)
	case "storagecertnamelist":
		response, err = h.storageCertNameList(request.Data)
	case "storagecertdata":
		response, err = h.storageCertData(request.Data)
	case "storagecertupload":
		response, err = h.uploadCertToStorage(request.Data)
	case "storagecertremove":
		err = h.removeCertFromStorage(request.Data)
	case "storagecertdownload":
		response, err = h.downloadCertFromStorage(request.Data)
	case "domainassign":
		response, err = h.assignCertificateToDomain(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func (h *Handler) issueCertificateToDomain(data interface{}) (*agentintegration.Certificate, error) {
	var certData agentintegration.CertificateIssueRequestData
	err := mapstructure.Decode(data, &certData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	return h.certificateManager.Issue(certData)
}

func (h *Handler) uploadCertificateToDomain(data interface{}) (*agentintegration.Certificate, error) {
	var requestData agentintegration.CertificateUploadRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	if requestData.ServerName == "" {
		return nil, errors.New("domain name is missed")
	}

	return h.certificateManager.Upload(requestData.ServerName, requestData.WebServer, requestData.PemCertificate)
}

func (h *Handler) storageCertNameList(data interface{}) (*agentintegration.StorageCertificateNameList, error) {
	certList, err := h.certificateManager.GetStorageCertList()
	if err != nil {
		return nil, err
	}
	certNameList := agentintegration.StorageCertificateNameList{
		CertNameList: certList,
	}
	return &certNameList, nil
}

func (h *Handler) storageCertData(data interface{}) (*agentintegration.Certificate, error) {
	certName, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid certificate name data is provided")
	}

	return h.certificateManager.GetStorageCertData(certName)
}

func (h *Handler) uploadCertToStorage(data interface{}) (*agentintegration.Certificate, error) {
	var requestData agentintegration.CertificateUploadRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	if requestData.CertName == "" {
		return nil, errors.New("certificate name is missed")
	}

	storage := GetDefaultCertStorage()
	_, err = storage.AddPemCertificate(requestData.CertName, requestData.PemCertificate)
	if err != nil {
		return nil, err
	}

	return storage.GetCertificate(requestData.CertName)
}

func (h *Handler) removeCertFromStorage(data interface{}) error {
	certName, ok := data.(string)
	if !ok {
		return errors.New("invalid certificate name data is provided")
	}

	storage := GetDefaultCertStorage()
	return storage.RemoveCertificate(certName)
}

func (h *Handler) downloadCertFromStorage(data interface{}) (*agentintegration.CertificateDownloadResponseData, error) {
	certName, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid certificate name data is provided")
	}

	storage := GetDefaultCertStorage()
	certPath, certContent, err := storage.GetCertificateAsString(certName)
	if err != nil {
		return nil, err
	}

	var certDownloadResponse agentintegration.CertificateDownloadResponseData
	certDownloadResponse.CertFileName = filepath.Base(certPath)
	certDownloadResponse.CertContent = certContent

	return &certDownloadResponse, nil
}

func (h *Handler) assignCertificateToDomain(data interface{}) (*agentintegration.Certificate, error) {
	var certData agentintegration.CertificateAssignRequestData
	err := mapstructure.Decode(data, &certData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	return h.certificateManager.Assign(certData)
}

func GetHandler(logger logger.LoggerInterface) (router.HandlerInterface, error) {
	certManager, err := GetCertificateManager(logger)
	if err != nil {
		return nil, err
	}

	return &Handler{logger: logger, certificateManager: certManager}, nil
}
