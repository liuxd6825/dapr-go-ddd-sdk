package mongo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dapr/durabletask-go/api"
	"github.com/dapr/durabletask-go/api/helpers"
	"github.com/dapr/durabletask-go/api/protos"
	"github.com/dapr/durabletask-go/backend"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/daos"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/daos/idao"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var emptyString string = ""

var errNoWorkItems = errors.New("no work items were found")

type Options struct {
	OrchestrationLockTimeout time.Duration
	ActivityLockTimeout      time.Duration
}

type mongoBackend struct {
	client       *mongo.Client
	options      *Options
	db           *ddd_mongodb.MongoDB
	logger       backend.Logger
	instancesDao *idao.Dao[*Instances]
	historyDao   *idao.Dao[*History]
	newEventsDao *idao.Dao[*NewEvents]
	newTasksDao  *idao.Dao[*NewTasks]
}

func (m *mongoBackend) CreateTaskHub(ctx context.Context) error {
	m.instancesDao = dao.NewDao[*Instances]("Instances")
	m.historyDao = mongo_dao.NewDao[*History]("History", &mongo_dao.RepositoryOptions{
		MongoDB: m.db,
	})
	m.newEventsDao = mongo_dao.NewDao[*NewEvents]("NewEvents", &mongo_dao.RepositoryOptions{
		MongoDB: m.db,
	})
	m.newTasksDao = mongo_dao.NewDao[*NewTasks]("NewTasks", &mongo_dao.RepositoryOptions{
		MongoDB: m.db,
	})
	return nil
}

func (m *mongoBackend) DeleteTaskHub(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) Start(ctx context.Context) error {
	return nil
}

func (m *mongoBackend) Stop(ctx context.Context) error {
	return nil
}

func (m *mongoBackend) CreateOrchestrationInstance(ctx context.Context, event *backend.HistoryEvent, options ...backend.OrchestrationIdReusePolicyOptions) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		var instanceID string
		if instanceID, err = m.createOrchestrationInstanceInternal(ctx, event); err != nil {
			return err
		}
		eventPayload, err := backend.MarshalHistoryEvent(event)
		if err != nil {
			return err
		}

		err = m.newEventsDao.Insert(
			ctx,
			&NewEvents{InstanceID: instanceID, EventPayload: eventPayload},
		)

		if err != nil {
			return fmt.Errorf("failed to insert row into [NewEvents] table: %w", err)
		}

		return session.CommitTransaction(ctx)
	})

	if err != nil {
		return err
	} else {
		fmt.Println("Transaction committed successfully")
	}

	return nil
}

func (m *mongoBackend) createOrchestrationInstanceInternal(ctx context.Context, e *backend.HistoryEvent, opts ...backend.OrchestrationIdReusePolicyOptions) (string, error) {
	if e == nil {
		return "", errors.New("HistoryEvent must be non-nil")
	} else if e.Timestamp == nil {
		return "", errors.New("HistoryEvent must have a non-nil timestamp")
	}

	startEvent := e.GetExecutionStarted()
	if startEvent == nil {
		return "", errors.New("HistoryEvent must be an ExecutionStartedEvent")
	}
	instanceID := startEvent.OrchestrationInstance.InstanceId

	policy := &protos.OrchestrationIdReusePolicy{}

	for _, opt := range opts {
		opt(policy)
	}

	rows, err := m.insertOrIgnoreInstanceTableInternal(ctx, e, startEvent)
	if err != nil {
		return "", err
	}

	// instance with same ID already exists
	if rows <= 0 {
		return instanceID, m.handleInstanceExists(ctx, startEvent, policy, e)
	}
	return instanceID, nil
}

func (m *mongoBackend) handleInstanceExists(ctx context.Context, startEvent *protos.ExecutionStartedEvent, policy *protos.OrchestrationIdReusePolicy, e *backend.HistoryEvent) error {
	// query RuntimeStatus for the existing instance
	runtimeStatus, _, err := m.instancesDao.FindById(ctx, "", startEvent.OrchestrationInstance.InstanceId)
	if errors.Is(err, mongo.ErrNilDocument) {
		return api.ErrInstanceNotFound
	} else if err != nil {
		return fmt.Errorf("failed to scan the Instances table result: %w", err)
	}

	// status not match, return instance duplicate error
	if !isStatusMatch(policy.OperationStatus, helpers.FromRuntimeStatusString(runtimeStatus.RuntimeStatus)) {
		return api.ErrDuplicateInstance
	}

	// status match
	switch policy.Action {
	case protos.CreateOrchestrationAction_IGNORE:
		// Log an warning message and ignore creating new instance
		m.logger.Warnf("An instance with ID '%s' already exists; dropping duplicate create request", startEvent.OrchestrationInstance.InstanceId)
		return api.ErrIgnoreInstance
	case protos.CreateOrchestrationAction_TERMINATE:
		// terminate existing instance
		if err := m.cleanupOrchestrationStateInternal(ctx, api.InstanceID(startEvent.OrchestrationInstance.InstanceId), false); err != nil {
			return fmt.Errorf("failed to cleanup orchestration status: %w", err)
		}
		// create a new instance
		var rows int64
		if rows, err = m.insertOrIgnoreInstanceTableInternal(ctx, e, startEvent); err != nil {
			return err
		}

		// should never happen, because we clean up instance before create new one
		if rows <= 0 {
			return fmt.Errorf("failed to insert into [Instances] table because entry already exists.")
		}
		return nil
	}
	// default behavior
	return api.ErrDuplicateInstance
}

func (m *mongoBackend) cleanupOrchestrationStateInternal(ctx context.Context, id api.InstanceID, requireCompleted bool) error {

	row, err := m.instancesDao.Count(ctx, "", string(id))
	if err != nil {
		return fmt.Errorf("failed to query for instance existence: %w", err)
	}

	if row != 1 {
		return fmt.Errorf("instance does not exist ")
	}

	if requireCompleted {
		// purge orchestration in ['COMPLETED', 'FAILED', 'TERMINATED']
		err = m.instancesDao.DeleteByFilter(ctx, "", fmt.Sprintf("[InstanceID] = ? AND [RuntimeStatus] IN ('COMPLETED', 'FAILED', 'TERMINATED')", string(id)))
		if err != nil {
			return fmt.Errorf("failed to delete from the Instances table: %w", err)
		}
	} else {
		// clean up orchestration in all [RuntimeStatus]
		err = m.instancesDao.DeleteById(ctx, "", string(id))
		if err != nil {
			return fmt.Errorf("failed to delete from the Instances table: %w", err)
		}
	}

	err = m.historyDao.Delete(ctx, &History{InstanceID: string(id)})
	if err != nil {
		return fmt.Errorf("failed to delete from History table: %w", err)
	}

	err = m.newEventsDao.Delete(ctx, &NewEvents{InstanceID: string(id)})
	if err != nil {
		return fmt.Errorf("failed to delete from NewEvents table: %w", err)
	}

	err = m.newTasksDao.Delete(ctx, &NewTasks{InstanceID: string(id)})
	if err != nil {
		return fmt.Errorf("failed to delete from NewTasks table: %w", err)
	}
	return nil
}

func (m *mongoBackend) insertOrIgnoreInstanceTableInternal(ctx context.Context, e *backend.HistoryEvent, startEvent *protos.ExecutionStartedEvent) (int64, error) {

	err := m.instancesDao.Insert(ctx,
		&Instances{
			InstanceID:    startEvent.OrchestrationInstance.InstanceId,
			Name:          startEvent.Name,
			Version:       startEvent.Version.GetValue(),
			ExecutionID:   startEvent.OrchestrationInstance.ExecutionId.GetValue(),
			Input:         startEvent.Input.GetValue(),
			RuntimeStatus: "PENDING",
			CreatedTime:   e.Timestamp.AsTime(),
		},
	)
	if err != nil {
		return -1, fmt.Errorf("failed to insert into [Instances] table: %w", err)
	}

	return 1, nil
}

func (m *mongoBackend) AddNewOrchestrationEvent(ctx context.Context, iid api.InstanceID, e *backend.HistoryEvent) error {
	if e == nil {
		return errors.New("HistoryEvent must be non-nil")
	} else if e.Timestamp == nil {
		return errors.New("HistoryEvent must have a non-nil timestamp")
	}

	eventPayload, err := backend.MarshalHistoryEvent(e)
	if err != nil {
		return err
	}

	err = m.newEventsDao.Insert(
		ctx,
		&NewEvents{
			InstanceID:   string(iid),
			EventPayload: eventPayload,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to insert row into [NewEvents] table: %w", err)
	}

	return nil
}

func (m *mongoBackend) GetOrchestrationWorkItem(ctx context.Context) (*backend.OrchestrationWorkItem, error) {
	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	now := time.Now().UTC()
	newLockExpiration := now.Add(m.options.OrchestrationLockTimeout)

	// Place a lock on an orchestration instance that has new events that are ready to be executed.

	/*
		row := tx.QueryRowContext(
			ctx,
			`UPDATE Instances SET [LockedBy] = ?, [LockExpiration] = ?
			WHERE [rowid] = (
				SELECT [rowid] FROM Instances I
				WHERE (I.[LockExpiration] IS NULL OR I.[LockExpiration] < ?) AND EXISTS (
					SELECT 1 FROM NewEvents E
					WHERE E.[InstanceID] = I.[InstanceID] AND (E.[VisibleTime] IS NULL OR E.[VisibleTime] < ?)
				)
				LIMIT 1
			) RETURNING [InstanceID]`,
			be.workerName,     // LockedBy for Instances table
			newLockExpiration, // Updated LockExpiration for Instances table
			now,               // LockExpiration for Instances table
			now,               // VisibleTime for NewEvents table
		)
	*/

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to query for orchestration work-items: %w", err)
	}

	var instanceID string
	if err := row.Scan(&instanceID); err != nil {
		if err == sql.ErrNoRows {
			// No new events to process
			return nil, errNoWorkItems
		}

		return nil, fmt.Errorf("failed to scan the orchestration work-item: %w", err)
	}

	// TODO: Get all the unprocessed events associated with the locked instance
	events, err := tx.QueryContext(
		ctx,
		`UPDATE NewEvents SET [DequeueCount] = [DequeueCount] + 1, [LockedBy] = ? WHERE rowid IN (
			SELECT rowid FROM NewEvents
			WHERE [InstanceID] = ? AND ([VisibleTime] IS NULL OR [VisibleTime] <= ?)
			LIMIT 1000
		)
		RETURNING [EventPayload], [DequeueCount]`,
		be.workerName,
		instanceID,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query for orchestration work-items: %w", err)
	}

	maxDequeueCount := int32(0)

	newEvents := make([]*protos.HistoryEvent, 0, 10)
	for events.Next() {
		var eventPayload []byte
		var dequeueCount int32
		if err := events.Scan(&eventPayload, &dequeueCount); err != nil {
			return nil, fmt.Errorf("failed to read history event: %w", err)
		}

		if dequeueCount > maxDequeueCount {
			maxDequeueCount = dequeueCount
		}

		e, err := backend.UnmarshalHistoryEvent(eventPayload)
		if err != nil {
			return nil, err
		}

		newEvents = append(newEvents, e)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to update orchestration work-item: %w", err)
	}

	wi := &backend.OrchestrationWorkItem{
		InstanceID: api.InstanceID(instanceID),
		NewEvents:  newEvents,
		LockedBy:   be.workerName,
		RetryCount: maxDequeueCount - 1,
	}

	return wi, nil
}

func (m *mongoBackend) GetOrchestrationRuntimeState(ctx context.Context, item *backend.OrchestrationWorkItem) (*backend.OrchestrationRuntimeState, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) GetOrchestrationMetadata(ctx context.Context, id api.InstanceID) (*backend.OrchestrationMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) CompleteOrchestrationWorkItem(ctx context.Context, item *backend.OrchestrationWorkItem) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) AbandonOrchestrationWorkItem(ctx context.Context, item *backend.OrchestrationWorkItem) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) GetActivityWorkItem(ctx context.Context) (*backend.ActivityWorkItem, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) CompleteActivityWorkItem(ctx context.Context, item *backend.ActivityWorkItem) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) AbandonActivityWorkItem(ctx context.Context, item *backend.ActivityWorkItem) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoBackend) PurgeOrchestrationState(ctx context.Context, id api.InstanceID) error {
	//TODO implement me
	panic("implement me")
}

func NewMongoBackend(dbClient *mongo.Client, logger backend.Logger) backend.Backend {
	return &mongoBackend{
		dbClient: dbClient,
		logger:   logger,
	}
}

func isStatusMatch(statuses []protos.OrchestrationStatus, runtimeStatus protos.OrchestrationStatus) bool {
	for _, status := range statuses {
		if status == runtimeStatus {
			return true
		}
	}
	return false
}
