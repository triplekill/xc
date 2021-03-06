package vulns

import (
	"fmt"
	"net"
	"strconv"

	"../shell"
	"../utils"
)

// https://www.catalog.update.microsoft.com/Home.aspx

// Vuln ...
type Vuln interface {
	Check(build string, version string, kbs []string) bool
	Name() string
	Description() string
}

// Check for common windows vulnerabilities
func Check(c net.Conn) {
	var raw string
	raw, _ = shell.ExecOut("ver")
	build := utils.GetBuild(raw)
	version, _ := shell.ExecPSOut(`(Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion").ReleaseId`) // e.g. 1909
	raw, _ = shell.ExecPSOut("get-hotfix")
	hotfixes := utils.GetHotfixes(raw)

	vulns := []Vuln{}
	vulns = append(vulns, &CVE20191315{})
	vulns = append(vulns, &CVE20200668{})
	vulns = append(vulns, &CVE20200787{})
	vulns = append(vulns, &CVE20200796{})
	vulns = append(vulns, &CVE20200863{})

	vulnerable := false
	result := ""
	for _, vuln := range vulns {
		if vuln.Check(build, version, hotfixes) {
			result += fmt.Sprintf("[+] vulnerable to %s\n", vuln.Name())
			result += fmt.Sprintf("    - %s\n", vuln.Description())
			vulnerable = true
		}
	}
	if !vulnerable {
		versionInt, err := strconv.Atoi(version)
		if err == nil {
			if versionInt < 1607 {
				result += "[-] OS not supported (<1607)"
			}
		}
		result += "[-] No common vulnerabilities found.."
	}
	result += "\n"
	// Debug
	c.Write([]byte(result))
}
