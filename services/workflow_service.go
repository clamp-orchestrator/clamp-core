package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"log"
)

func SaveWorkflow(workflowReq models.Workflow) (models.Workflow, error) {
	log.Printf("Saving worflow %v", workflowReq)
	workflow, err := repository.DB.SaveWorkflow(workflowReq)
	if err != nil {
		log.Printf("Failed to save workflow: %v, error: %s\n", workflow, err.Error())
	} else {
		log.Printf("Saved worflow %v", workflow)
	}
	return workflow, err
}

func FindWorkflowByName(workflowName string) (models.Workflow, error) {
	log.Printf("Finding workflow by name : %s", workflowName)
	workflow, err := repository.GetDB().FindWorkflowByName(workflowName)
	if err != nil {
		log.Printf("No record found with given workflow name %s, error: %s\n", workflowName, err.Error())
	}
	return workflow, err
}

func GetWorkflows(pageNumber int, pageSize int, filters map[string]string) ([]models.Workflow, error) {
	log.Printf("Getting workflows for pageNumber: %d, pageSize: %d", pageNumber, pageSize)
	workflows, err := repository.GetDB().GetWorkflows(pageNumber, pageSize, filters)
	if err != nil {
		log.Printf("Failed to fetch worflows for pageNumber: %d, pageSize: %d", pageNumber, pageSize)
	}
	return workflows, err
}
