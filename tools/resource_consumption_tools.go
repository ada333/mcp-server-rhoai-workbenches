package tools

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func parseResourceValue(s string) (float64, string, error) {
	if s == "" || s == "0" {
		return 0, "", nil
	}

	re := regexp.MustCompile(`^([\d.]+)(.*)$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, "", fmt.Errorf("invalid resource format: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse number from %s: %w", s, err)
	}

	unit := strings.TrimSpace(matches[2])
	return value, unit, nil
}

// sumResourceValues sums multiple resource strings and returns the result with the (first) unit
func sumResourceValues(values []string) string {
	if len(values) == 0 {
		return "0"
	}

	var total float64
	var unit string

	for _, v := range values {
		if v == "" || v == "0" {
			continue
		}

		val, u, err := parseResourceValue(v)
		if err != nil {
			continue
		}

		total += val
		// if there would be multiple units, it wouldnt work correctly
		if unit == "" && u != "" {
			unit = u
		}
	}
	if total == 0 {
		return "0"
	}
	if total == float64(int64(total)) {
		return fmt.Sprintf("%d%s", int64(total), unit)
	}
	return fmt.Sprintf("%.2f%s", total, unit)
}

func ListResourceConsumptionPerWorkbench(ctx context.Context, req *mcp.CallToolRequest, input core.ListResourceConsumptionPerWorkbenchInput) (*mcp.CallToolResult, core.ListResourceConsumptionOutput, error) {
	_, workbenches, err := ListWorkbenches(ctx, req, core.ListWorkbenchesInput{Namespace: input.Namespace})
	if err != nil {
		return nil, core.ListResourceConsumptionOutput{}, err
	}
	// theoretically this is uneffective - there could be function that returns the workbench by name and namespace
	// and not list all workbenches (linear complexity)
	for _, wb := range workbenches.Workbenches {
		if wb.Name == input.WorkbenchName {
			return nil, core.ListResourceConsumptionOutput{
				CPUUsage:    wb.CPUUsage,
				MemoryUsage: wb.MemoryUsage,
				DiskUsage:   wb.DiskUsage,
			}, nil
		}
	}

	return nil, core.ListResourceConsumptionOutput{}, fmt.Errorf("workbench %s not found in namespace %s", input.WorkbenchName, input.Namespace)
}

func ListResourceConsumptionPerNamespace(ctx context.Context, req *mcp.CallToolRequest, input core.ListResourceConsumptionPerNamespaceInput) (*mcp.CallToolResult, core.ListResourceConsumptionOutput, error) {
	_, workbenches, err := ListWorkbenches(ctx, req, core.ListWorkbenchesInput(input))
	if err != nil {
		return nil, core.ListResourceConsumptionOutput{}, err
	}

	var cpuValues, memoryValues, diskValues []string
	for _, wb := range workbenches.Workbenches {
		cpuValues = append(cpuValues, wb.CPUUsage)
		memoryValues = append(memoryValues, wb.MemoryUsage)
		diskValues = append(diskValues, wb.DiskUsage)
	}

	return nil, core.ListResourceConsumptionOutput{
		CPUUsage:    sumResourceValues(cpuValues),
		MemoryUsage: sumResourceValues(memoryValues),
		DiskUsage:   sumResourceValues(diskValues),
	}, nil
}

func ListResourceConsumptionPerUser(ctx context.Context, req *mcp.CallToolRequest, input core.ListResourceConsumptionPerUserInput) (*mcp.CallToolResult, core.ListResourceConsumptionOutput, error) {
	_, workbenches, err := ListAllWorkbenches(ctx, req, core.ListWorkbenchesInput{})
	if err != nil {
		return nil, core.ListResourceConsumptionOutput{}, err
	}

	var cpuValues, memoryValues, diskValues []string
	for _, wb := range workbenches.Workbenches {
		if wb.User == input.User {
			cpuValues = append(cpuValues, wb.CPUUsage)
			memoryValues = append(memoryValues, wb.MemoryUsage)
			diskValues = append(diskValues, wb.DiskUsage)
		}
	}

	return nil, core.ListResourceConsumptionOutput{
		CPUUsage:    sumResourceValues(cpuValues),
		MemoryUsage: sumResourceValues(memoryValues),
		DiskUsage:   sumResourceValues(diskValues),
	}, nil
}

func ListResourceConsumptionPerCluster(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, core.ListResourceConsumptionOutput, error) {
	_, workbenches, err := ListAllWorkbenches(ctx, req, core.ListWorkbenchesInput{})
	if err != nil {
		return nil, core.ListResourceConsumptionOutput{}, err
	}

	var cpuValues, memoryValues, diskValues []string
	for _, wb := range workbenches.Workbenches {
		cpuValues = append(cpuValues, wb.CPUUsage)
		memoryValues = append(memoryValues, wb.MemoryUsage)
		diskValues = append(diskValues, wb.DiskUsage)
	}

	return nil, core.ListResourceConsumptionOutput{
		CPUUsage:    sumResourceValues(cpuValues),
		MemoryUsage: sumResourceValues(memoryValues),
		DiskUsage:   sumResourceValues(diskValues),
	}, nil
}
