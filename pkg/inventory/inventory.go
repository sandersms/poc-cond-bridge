// SPDX-License-Identifier: Apache-2.0

// Package inventory is the base inventory functions
package inventory

import (
	"context"
	"fmt"
	"log"

	ipb "github.com/opiproject/opi-api/inventory/v1/gen/go"

	"github.com/jaypipes/ghw"
)

// Server contains inventory related OPI services
type Server struct {
	ipb.UnimplementedInventoryServiceServer
}

// GetInventory returns inventory information
func (s *Server) GetInventory(_ context.Context, in *ipb.GetInventoryRequest) (*ipb.Inventory, error) {
	log.Printf("GetInventory: Received from client: %v", in)

	cpu, err := ghw.CPU()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v", err)
		// return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	fmt.Printf("%v\n", cpu)

	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}
	fmt.Println(memory.String())

	chassis, err := ghw.Chassis()
	if err != nil {
		fmt.Printf("Error getting chassis info: %v", err)
	}
	fmt.Printf("%v\n", chassis)

	bios, err := ghw.BIOS()
	if err != nil {
		fmt.Printf("Error getting BIOS info: %v", err)
	}
	fmt.Printf("%v\n", bios)

	baseboard, err := ghw.Baseboard()
	if err != nil {
		fmt.Printf("Error getting baseboard info: %v", err)
	}
	fmt.Printf("%v\n", baseboard)

	product, err := ghw.Product()
	if err != nil {
		fmt.Printf("Error getting product info: %v", err)
	}
	fmt.Printf("%v\n", product)

	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting pci info: %v", err)
	}
	Blobarray := make([]*ipb.PCIeDeviceInfo, len(pci.Devices))
	for i, r := range pci.Devices {
		fmt.Printf("PCI=%v\n", r)
		Blobarray[i] = &ipb.PCIeDeviceInfo{Driver: r.Driver, Address: r.Address, Vendor: r.Vendor.Name, Product: r.Product.Name, Revision: r.Revision, Subsystem: r.Subsystem.Name, Class: r.Class.Name, Subclass: r.Subclass.Name}
	}

	return &ipb.Inventory{
		Bios:      &ipb.BIOSInfo{Vendor: bios.Vendor, Version: bios.Version, Date: bios.Date},
		System:    &ipb.SystemInfo{Family: product.Family, Name: product.Name, Vendor: product.Vendor, SerialNumber: product.SerialNumber, Uuid: product.UUID, Sku: product.SKU, Version: product.Version},
		Baseboard: &ipb.BaseboardInfo{AssetTag: baseboard.AssetTag, SerialNumber: baseboard.SerialNumber, Vendor: baseboard.Vendor, Version: baseboard.Version, Product: baseboard.Product},
		Chassis:   &ipb.ChassisInfo{AssetTag: chassis.AssetTag, SerialNumber: chassis.SerialNumber, Type: chassis.Type, TypeDescription: chassis.TypeDescription, Vendor: chassis.Vendor, Version: chassis.Version},
		Processor: &ipb.CPUInfo{TotalCores: int32(cpu.TotalCores), TotalThreads: int32(cpu.TotalThreads)},
		Memory:    &ipb.MemoryInfo{TotalPhysicalBytes: memory.TotalPhysicalBytes, TotalUsableBytes: memory.TotalUsableBytes},
		Pci:       Blobarray,
	}, nil
}