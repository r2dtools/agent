package certificates

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates/commondir"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/router"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/agentintegration"
)

type Handler struct {
	certificateManager *CertificateManager
	logger             logger.Logger
	config             *config.Config
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
		response, err = h.storageCertNameList()
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
	case "commondirstatus":
		response, err = h.commonDirStatus(request.Data)
	case "changecommondirstatus":
		err = h.changeCommonDirStatus(request.Data)
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

func (h *Handler) storageCertNameList() (*agentintegration.StorageCertificateNameList, error) {
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

	storage := GetDefaultCertStorage(h.config)
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

	storage := GetDefaultCertStorage(h.config)
	return storage.RemoveCertificate(certName)
}

func (h *Handler) downloadCertFromStorage(data interface{}) (*agentintegration.CertificateDownloadResponseData, error) {
	certName, ok := data.(string)
	if !ok {
		return nil, errors.New("invalid certificate name data is provided")
	}

	storage := GetDefaultCertStorage(h.config)
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

func (h *Handler) commonDirStatus(data interface{}) (*agentintegration.CommonDirStatusResponseData, error) {
	var requestData agentintegration.CommonDirChangeStatusRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return nil, fmt.Errorf("invalid common dir status request data: %v", err)
	}

	options := h.config.ToMap()
	wServer, err := webserver.GetWebServer(requestData.WebServer, options)

	if err != nil {
		return nil, err
	}

	webServerReverter := &reverter.Reverter{
		HostMng: wServer.GetVhostManager(),
		Logger:  h.logger,
	}

	commonDirManager, err := commondir.GetCommonDirManager(wServer, webServerReverter, h.logger, options)

	if err != nil {
		return nil, err
	}

	status := commonDirManager.IsCommonDirEnabled(requestData.ServerName)

	return &agentintegration.CommonDirStatusResponseData{Status: status}, nil
}

func (h *Handler) changeCommonDirStatus(data interface{}) error {
	var requestData agentintegration.CommonDirChangeStatusRequestData
	err := mapstructure.Decode(data, &requestData)

	if err != nil {
		return fmt.Errorf("invalid common dir status request data: %v", err)
	}

	options := h.config.ToMap()
	wServer, err := webserver.GetWebServer(requestData.WebServer, options)

	if err != nil {
		return err
	}

	webServerReverter := &reverter.Reverter{
		HostMng: wServer.GetVhostManager(),
		Logger:  h.logger,
	}

	commonDirManager, err := commondir.GetCommonDirManager(wServer, webServerReverter, h.logger, options)

	if err != nil {
		return err
	}

	if requestData.Status {
		err = commonDirManager.EnableCommonDir(requestData.ServerName)
	} else {
		err = commonDirManager.DisableCommonDir(requestData.ServerName)
	}

	return err
}

func GetHandler(config *config.Config, logger logger.Logger) (router.HandlerInterface, error) {
	certManager, err := GetCertificateManager(config, logger)

	if err != nil {
		return nil, err
	}

	return &Handler{
		logger:             logger,
		certificateManager: certManager,
		config:             config,
	}, nil
}
