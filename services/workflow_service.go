package services

import (
	"clamp-core/models"
	"clamp-core/repository"

	log "github.com/sirupsen/logrus"
)

func SaveWorkflow(workflowReq models.Workflow) (models.Workflow, error) {
	log.Debugf("Saving worflow %v", workflowReq)
	workflow, err := repository.GetDB().SaveWorkflow(workflowReq)
	if err != nil {
		log.Errorf("Failed to save workflow: %v, error: %s", workflow, err.Error())
	} else {
		log.Debugf("Saved worflow %v", workflow)
	}
	return workflow, err
}

func FindWorkflowByName(workflowName string) (models.Workflow, error) {
	log.Debugf("Finding workflow by name : %s", workflowName)
	workflow, err := repository.GetDB().FindWorkflowByName(workflowName)
	if err != nil {
		log.Errorf("No record found with given workflow name %s, error: %s", workflowName, err.Error())
	}
	return workflow, err
}

// DeleteWorkflow will delete the existing workflow by name
// This method is not exposed an an API. It is implemented for running a test scenario.
func DeleteWorkflowByName(workflowName string) error {
	log.Debugf("Deleting workflow by name : %s", workflowName)
	err := repository.GetDB().DeleteWorkflowByName(workflowName)
	if err != nil {
		log.Debugf("No record found with given workflow name %s, error: %s", workflowName, err.Error())
	}
	return err
}

// GetWorkflows is used to fetch all the workflows for the GET call API
// Implements a pagination approach
// Also supports filters
func GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.Workflow, int, error) {
	log.Debugf("Getting workflows for pageNumber: %d, pageSize: %d", pageNumber, pageSize)
	workflows, totalWorkflowsCount, err := repository.GetDB().GetWorkflows(pageNumber, pageSize, sortBy)
	if err != nil {
		log.Debugf("Failed to fetch worflows for pageNumber: %d, pageSize: %d, sortBy %v", pageNumber, pageSize, sortBy)
	}
	return workflows, totalWorkflowsCount, err
}
