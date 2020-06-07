package handlers

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"squzy/internal/helpers"
	"time"
)

type Handlers interface {
	GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error)
	GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error)
	GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error)
	GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error)
	GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error)
	GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error)
	RunScheduler(ctx context.Context, id string) error
	StopScheduler(ctx context.Context, id string) error
	RemoveScheduler(ctx context.Context, id string) error
	AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error)
	RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error)
	SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error)
	GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error)
	GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error)
	GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error)
	GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error)
	GetApplicationList(ctx context.Context) ([]*apiPb.Application, error)
	GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error)
}

const (
	defaultRequestTimeout = time.Second * 30
)

type handlers struct {
	agentClient                 apiPb.AgentServerClient
	monitoringClient            apiPb.SchedulersExecutorClient
	storageClient               apiPb.StorageClient
	applicationMonitoringClient apiPb.ApplicationMonitoringClient
}

func (h *handlers) SaveTransaction(ctx context.Context, rq *apiPb.TransactionInfo) (*empty.Empty, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	r, err := h.applicationMonitoringClient.SaveTransaction(c, rq)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *handlers) GetSchedulerHistoryByID(ctx context.Context, rq *apiPb.GetSchedulerInformationRequest) (*apiPb.GetSchedulerInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	r, err := h.storageClient.GetSchedulerInformation(c, rq)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *handlers) GetAgentHistoryByID(ctx context.Context, rq *apiPb.GetAgentInformationRequest) (*apiPb.GetAgentInformationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	r, err := h.storageClient.GetAgentInformation(c, rq)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (h *handlers) RunScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Run(c, &apiPb.RunRequest{
		Id: id,
	})
	return err
}

func (h *handlers) AddScheduler(ctx context.Context, scheduler *apiPb.AddRequest) (*apiPb.AddResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.monitoringClient.Add(c, scheduler)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *handlers) StopScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Stop(c, &apiPb.StopRequest{
		Id: id,
	})
	return err
}

func (h *handlers) RemoveScheduler(ctx context.Context, id string) error {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	_, err := h.monitoringClient.Remove(c, &apiPb.RemoveRequest{
		Id: id,
	})
	return err
}

func (h *handlers) GetSchedulerList(ctx context.Context) ([]*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	list, err := h.monitoringClient.GetSchedulerList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Lists, nil
}

func (h *handlers) GetSchedulerByID(ctx context.Context, id string) (*apiPb.Scheduler, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	scheduler, err := h.monitoringClient.GetSchedulerById(c, &apiPb.GetSchedulerByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return scheduler, nil
}

func (h *handlers) GetAgentList(ctx context.Context) ([]*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	list, err := h.agentClient.GetAgentList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return list.Agents, nil
}

func (h *handlers) GetAgentByID(ctx context.Context, id string) (*apiPb.AgentItem, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	agent, err := h.agentClient.GetAgentById(c, &apiPb.GetAgentByIdRequest{
		AgentId: id,
	})
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (h *handlers) RegisterApplication(ctx context.Context, rq *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	app, err := h.applicationMonitoringClient.InitializeApplication(c, rq)
	if err != nil {
		return nil, err
	}
	return app, err
}

func (h *handlers) GetSchedulerUptime(ctx context.Context, rq *apiPb.GetSchedulerUptimeRequest) (*apiPb.GetSchedulerUptimeResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.storageClient.GetSchedulerUptime(c, rq)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (h *handlers) GetTransactionGroups(ctx context.Context, req *apiPb.GetTransactionGroupRequest) (*apiPb.GetTransactionGroupResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.storageClient.GetTransactionsGroup(c, req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (h *handlers) GetTransactionsList(ctx context.Context, req *apiPb.GetTransactionsRequest) (*apiPb.GetTransactionsResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.storageClient.GetTransactions(c, req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (h *handlers) GetTransactionById(ctx context.Context, id string) (*apiPb.GetTransactionByIdResponse, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.storageClient.GetTransactionById(c, &apiPb.GetTransactionByIdRequest{
		TransactionId: id,
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (h *handlers) GetApplicationById(ctx context.Context, id string) (*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.applicationMonitoringClient.GetApplicationById(c, &apiPb.ApplicationByIdReuqest{
		ApplicationId: id,
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (h *handlers) GetApplicationList(ctx context.Context) ([]*apiPb.Application, error) {
	c, cancel := helpers.TimeoutContext(ctx, defaultRequestTimeout)
	defer cancel()
	res, err := h.applicationMonitoringClient.GetApplicationList(c, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return res.Applications, err
}

func New(
	agentClient apiPb.AgentServerClient,
	monitoringClient apiPb.SchedulersExecutorClient,
	storageClient apiPb.StorageClient,
	applicationMonitoringClient apiPb.ApplicationMonitoringClient,
) Handlers {
	return &handlers{
		agentClient:                 agentClient,
		monitoringClient:            monitoringClient,
		storageClient:               storageClient,
		applicationMonitoringClient: applicationMonitoringClient,
	}
}
