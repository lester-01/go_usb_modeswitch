package main

import (
        "fmt"
        "log"
        "os"
        "os/exec"
        "strings"
        "time"
)

func main() {
        // Check if usb_modeswitch is installed.
        if _, err := exec.LookPath("usb_modeswitch"); err != nil {
                fmt.Println("usb_modeswitch not found. Installing...")

                // Run "sudo apt update"
                cmdUpdate := exec.Command("sudo", "apt", "update")
                // Redirect command output to the standard output/error.
                cmdUpdate.Stdout = os.Stdout
                cmdUpdate.Stderr = os.Stderr
                if err := cmdUpdate.Run(); err != nil {
                        log.Fatalf("Failed to update apt repositories: %v", err)
                }

                // Run "sudo apt install -y usb-modeswitch"
                cmdInstall := exec.Command("sudo", "apt", "install", "-y", "usb-modeswitch")
                cmdInstall.Stdout = os.Stdout
                cmdInstall.Stderr = os.Stderr
                if err := cmdInstall.Run(); err != nil {
                        log.Fatalf("Failed to install usb-modeswitch: %v", err)
                }
        }

        // List of Huawei mass storage mode IDs.
        storageModeIDs := []string{
                "12d1:1446", "12d1:1506", "12d1:151d", "12d1:1520", "12d1:1f01",
                "12d1:1f11", "12d1:1f16", "12d1:1f17", "12d1:1f19", "12d1:1f21",
                "12d1:1f22", "12d1:1f23", "12d1:1f25", "12d1:1f28", "12d1:1f29",
        }

        // Infinite loop to continuously check for Huawei dongles in storage mode.
        for {
                // Run the "lsusb" command.
                lsusbOutput, err := exec.Command("lsusb").Output()
                if err != nil {
                        log.Printf("Error running lsusb: %v", err)
                        time.Sleep(1 * time.Second)
                        continue
                }

                outputStr := string(lsusbOutput)
                // Loop through IDs and check if any are present in the lsusb output.
                for _, id := range storageModeIDs {
                        if strings.Contains(outputStr, id) {
                                // Split the ID into vendor and product parts.
                                parts := strings.Split(id, ":")
                                if len(parts) != 2 {
                                        log.Printf("Invalid ID format: %s", id)
                                        continue
                                }
                                vendor := parts[0]
                                product := parts[1]

                                fmt.Printf("Huawei dongle detected in storage mode (%s). Switching...\n", id)
                                // Switch the mode using usb_modeswitch.
                                cmdSwitch := exec.Command("sudo", "usb_modeswitch", "-J", "-v", vendor, "-p", product)
                                cmdSwitch.Stdout = os.Stdout
                                cmdSwitch.Stderr = os.Stderr
                                if err := cmdSwitch.Run(); err != nil {
                                        log.Printf("Failed to switch dongle mode for %s: %v", id, err)
                                }

                                // Allow time for the switch to complete.
                                //time.Sleep(5 * time.Second)
                        }
                }
                // Wait a bit before the next check.
                time.Sleep(1 * time.Second)
        }
}
