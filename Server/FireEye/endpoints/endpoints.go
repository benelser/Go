package EndPoints

import (
    "fmt"
	"io"
	"bytes"
	"net/http"
	"crypto/tls"
	b64 "encoding/base64"
	"errors"
	"encoding/json"
	"time"
	"net/url"
	"io/ioutil"
	. "../../Utils"
)

func Connect(server *Server) (e error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	var r string
	fmt.Println("Starting the application...")
	r = fmt.Sprintf("https://%s:%s/hx/api/v3/", server.IP, server.Port)
	server.RootURI = r
	tokenUri := fmt.Sprintf("%stoken", r)
	request, _ := http.NewRequest("GET", tokenUri, nil)
	request.Header.Set("Content-Type", "application/json")
	credCombo := fmt.Sprintf("%s:%s", server.Username, server.Password)
	creds := fmt.Sprintf("Basic %s",  b64.StdEncoding.EncodeToString([]byte(credCombo)))
	request.Header.Set("Authorization", creds)
	client := &http.Client{Transport :tr}
	response, err := client.Do(request)
	if err != nil {
		errMessage := fmt.Sprintf("The HTTP request failed with error %s\n", err)
		myErr := errors.New(errMessage)
		return myErr
	} else {

		// Get Token To use later
		// data, _ := ioutil.ReadAll(response.Body)
		server.Token = response.Header["X-Feapi-Token"][0]
	}

	return nil
}

func GetVersion(server *Server) (v version, e error) {
	res, _ := APICall(server, nil, "version", "GET")
	defer res.Body.Close()
	ver := version{}
	err := json.NewDecoder(res.Body).Decode(&ver)
	if err != nil{
		errMessage := fmt.Sprintf("Failed to decode JSON Version with: %s\n", err)
		myErr := errors.New(errMessage)
		return ver, myErr
	}
	return ver, nil
}

func NewIndicator(server *Server, name, description string) (success int) {
	// Indicators are created first, and then conditions are added to them. If a new condition
	// matches one that was created previously, it will be assigned the same condition_id.
	values := map[string]string{"create_text": name, "description": description}
	jsonValue, _ := json.Marshal(values)
	body := bytes.NewBuffer(jsonValue)
	res, _ := APICall(server, body, "indicators/Custom", "POST")
	if res.StatusCode != 201 {

		return 0
	}

	return 1
}

func GetLiveHosts(server *Server) (livehost []string) {

	// Call host endpoint with limit tof 60k
	res, err := APICall(server, nil, "hosts?limit=60000", "GET")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer res.Body.Close()
	// Decode reponse to static struct generated using:
	// https://mholt.github.io/json-to-go/
	listOfHostsresponse := ListOfHostResponse{}
	json.NewDecoder(res.Body).Decode(&listOfHostsresponse)

	// inititalize slice with len of 0
	liveHosts := make([]string, 0)
	// All hosts are nested in the Entries Key/Prop
	for _, computer := range listOfHostsresponse.Data.Entries {
		// Build Hosts Slice
		liveHosts = append(liveHosts, computer.Hostname)
	}

	// Filter out dups
	return SortUniqueStringSlice(liveHosts)
	
}

func GetHostIDByHostName(server *Server, hostname string) (hostID string, e error) {
	fmtURI := fmt.Sprintf("hosts?hostname=%s&limit=1", url.QueryEscape(hostname))
	res, _ := APICall(server, nil, fmtURI, "GET")
	defer res.Body.Close()
	listOfHostsresponse := ListOfHostResponse{}
	json.NewDecoder(res.Body).Decode(&listOfHostsresponse)
	if len(listOfHostsresponse.Data.Entries) == 0 {
		errMessage := fmt.Sprintf("%s", hostname)
		myErr := errors.New(errMessage)
		return " ", myErr
	}
	return listOfHostsresponse.Data.Entries[0].ID, nil
}

func APICall(server *Server, body io.Reader, uri, method string) (res *http.Response, e error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}

	fmtURI := fmt.Sprintf("%s%s", server.RootURI, uri)

	var request *http.Request
	if body != nil {
		req, _ := http.NewRequest(method, fmtURI, body)
		request = req
	} else if body == nil {
		req, _ := http.NewRequest(method, fmtURI, nil)
		request = req
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-FeApi-Token", server.Token)
	client := &http.Client{Transport :tr}
	response, err := client.Do(request)

	if err != nil {
		errMessage := fmt.Sprintf("The HTTP request failed with error %s\n", err)
		myErr := errors.New(errMessage)
		return nil, myErr
	}
	return response, nil
}

func GetHostSetIDByName(server *Server, hostSetName string) (hostSetID int, e error) {
	fmtURI := fmt.Sprintf("host_sets?name=%s&limit=1", url.QueryEscape(hostSetName))
	res, _ := APICall(server, nil, fmtURI, "GET")
	defer res.Body.Close()
	if res.StatusCode != 200 {

		errMessage := fmt.Sprintf("Failed to retrieve ID for: %s\n", hostSetName)
		myErr := errors.New(errMessage)
		return 0, myErr
		
	}
	hostSets := HostSet{}
	json.NewDecoder(res.Body).Decode(&hostSets)
	
	if len(hostSets.Data.Entries) == 0 {
		errMessage := fmt.Sprintf("Failed to retrieve ID for: %s\n", hostSetName)
		myErr := errors.New(errMessage)
		return 0, myErr
	}
	return hostSets.Data.Entries[0].ID, nil
}

func GetStaticHostSetHosts(server *Server, hostSetID int) (hosts []Host, hostnames []string){

	var existingHosts []Host
	var existingHostsNames []string
	// Create URI
	fmtURI := fmt.Sprintf("host_sets/%v/hosts", hostSetID)

	// Make the call
	res, _ := APICall(server, nil, fmtURI, "GET")
	defer res.Body.Close()
	listOfHostsresponse := ListOfHostResponse{}
	json.NewDecoder(res.Body).Decode(&listOfHostsresponse)

	for _, host := range listOfHostsresponse.Data.Entries {
		
		existingHosts = append(existingHosts, host)
		existingHostsNames = append(existingHostsNames, host.Hostname)
	}
	return existingHosts, existingHostsNames
}

func UpdateStaticHostSet(server *Server, computers []string, hostSetID int, hostSetName string) (hostsNotInFireEye []string, e error) {

	// define var to hold hosts ID's and failures
	var hostsByIDToAdd []string
	var hostsNotFound []string

	_, hostnames := GetStaticHostSetHosts(server, hostSetID)
	hostNamesMap := SliceToIntMap(hostnames)

	// Create URI
	fmtURI := fmt.Sprintf("host_sets/static/%v", hostSetID)

	// Loop though computers appending to Hosts ID slice
	for _, computer := range computers {
		if _, ok := hostNamesMap[computer]; ok {
			continue
		}
		id, err := GetHostIDByHostName(server, computer)
		if err != nil{
			hostsNotFound = append(hostsNotFound, err.Error())
			continue
		}

		hostsByIDToAdd = append(hostsByIDToAdd, id)
	}

	// Build PUT Request body this is nested struct object
	var changes []StaticHostsChanges
	change := StaticHostsChanges{
		Command : "change",
		Add: hostsByIDToAdd,
		Remove: []string{},
	}
	changes = append(changes, change)

	// Main struct object
	var request StaticHostSetBody
	request.Name = hostSetName

	// Assign nested struct var
	request.Changes = changes

	// See what it look likes
	fmt.Printf("%v\r\n", request)

	// Convert to JSON byte array for PUT body
	jsonValue, _ := json.Marshal(request)
	body := bytes.NewBuffer(jsonValue)

	// See what it look likes
	fmt.Printf("%v\r\n", string(jsonValue))
	
	// Make the call
	res, _ := APICall(server, body, fmtURI, "PUT")
	if res.StatusCode != 200 {
		errMessage := fmt.Sprintf("Failed to update HostSet %v\nWith StatusCode %v\r\n", hostSetName, res.StatusCode)
		myErr := errors.New(errMessage)
		return nil, myErr
	}

	return hostsNotFound, nil
}

func GetPolicyByName(server *Server, policyName string) (p Policy, e error){
	// Create URI
	fmtURI := fmt.Sprintf("policies?name=%s", url.QueryEscape(policyName))

	// Make the call
	res, _ := APICall(server, nil, fmtURI, "GET")
	defer res.Body.Close()
	listOfPolicyresponse := ListOfPolicyResponse{}
	json.NewDecoder(res.Body).Decode(&listOfPolicyresponse)
	if len(listOfPolicyresponse.Data.Entries) == 0 {
		errMessage := fmt.Sprintf("%s does not exists\n", policyName)
		myErr := errors.New(errMessage)
		return Policy{}, myErr
	}
	return listOfPolicyresponse.Data.Entries[0], nil
}

func PrintResponseBody(res *http.Response){
	// Read Body to get schema
	bodyBytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("%v", err)
    }
	bodyString := string(bodyBytes)
	fmt.Printf("Response: \n\n%v", bodyString)
}

func GetPolicyByID(server *Server, ID string) (p PolicyUpdate, e error)  {
	fmtURI := fmt.Sprintf("policies/%v", ID)
	
	// Make the call
	res, _ := APICall(server, nil, fmtURI, "GET")
	defer res.Body.Close()

	// // This object only has policy specific props that can be updated diffs slightly from the Policy Struct
	var response PolicyUpdate
	json.NewDecoder(res.Body).Decode(&response)

	if len(response.Data) == 0 {
		errMessage := fmt.Sprintf("%s does not exists or could not be accessed.\n", ID)
		myErr := errors.New(errMessage)
		return PolicyUpdate{}, myErr
		
	}
	return response, nil
}

func UpdatePolicy(server *Server, updatedObject *PolicyUpdate, policyObject Policy) (e error){
	// Create URI
	fmtURI := fmt.Sprintf("policies/%s", policyObject.ID)

	// Convert to JSON byte array for PUT body
	jsonValue, _ := json.Marshal(updatedObject)
	body := bytes.NewBuffer(jsonValue)

	// See what it look likes
	fmt.Printf("%v\r\n", string(jsonValue))

	// Make the call
	res, _ := APICall(server, body, fmtURI, "PUT")
	if res.StatusCode != 200 {
		fmt.Println("something went wrong")
		fmt.Printf("\n%v", res)
		errMessage := fmt.Sprintf("Failed to update Policy %v\nWith StatusCode %v\r\n", policyObject.Name, res.StatusCode)
		myErr := errors.New(errMessage)
		return myErr
	}

	return nil
}

type Server struct {
	Token string
	RootURI string
	IP string 
	Port string
	Username string
	Password string
}

type version struct {
	// Must be caps
	// One to one static map to API's JSON response
	Details []string 
	Route string
	Data struct{
		Version string
		IntelVersion string
		MsoVersion string
		ApplianceId string
		IntelLastUpdateTime string
		IsUpgraded bool
	}
	Message string
}

type ListOfHostResponse struct {
	Data struct {
		Total int `json:"total"`
		Query struct {
		} `json:"query"`
		Sort struct {
		} `json:"sort"`
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Entries []Host `json:"entries"`
	} `json:"data"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
	Route   string        `json:"route"`
}

type HostSet struct {
	Data struct {
		Total int `json:"total"`
		Query struct {
		} `json:"query"`
		Sort struct {
		} `json:"sort"`
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Entries []struct {
			ID       int    `json:"_id"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			Revision string `json:"_revision"`
			URL      string `json:"url"`
		} `json:"entries"`
	} `json:"data"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
	Route   string        `json:"route"`
}

type StaticHostSetBody struct {
	Name    string `json:"name"`
	Changes []StaticHostsChanges `json:"changes"`
	
}

type StaticHostsChanges struct{
	Command string   `json:"command"`
	Add     []string `json:"add"`
	Remove  []string `json:"remove"`
}

type Host struct {
	ID                         string `json:"_id"`
	AgentVersion               string `json:"agent_version"`
	ExcludedFromContainment    bool   `json:"excluded_from_containment"`
	ContainmentMissingSoftware bool   `json:"containment_missing_software"`
	ContainmentQueued          bool   `json:"containment_queued"`
	ContainmentState           string `json:"containment_state"`
	Stats                      struct {
		Acqs               int `json:"acqs"`
		AlertingConditions int `json:"alerting_conditions"`
		Alerts             int `json:"alerts"`
		ExploitAlerts      int `json:"exploit_alerts"`
		ExploitBlocks      int `json:"exploit_blocks"`
		MalwareAlerts      int `json:"malware_alerts"`
	} `json:"stats"`
	Hostname                  string      `json:"hostname"`
	Domain                    string      `json:"domain"`
	GmtOffsetSeconds          int         `json:"gmt_offset_seconds"`
	Timezone                  string      `json:"timezone"`
	PrimaryIPAddress          string      `json:"primary_ip_address"`
	LastAuditTimestamp        time.Time   `json:"last_audit_timestamp"`
	LastPollTimestamp         interface{} `json:"last_poll_timestamp"`
	LastPollIP                interface{} `json:"last_poll_ip"`
	ReportedClone             bool        `json:"reported_clone"`
	InitialAgentCheckin       time.Time   `json:"initial_agent_checkin"`
	URL                       string      `json:"url"`
	LastAlert                 interface{} `json:"last_alert"`
	LastExploitBlock          interface{} `json:"last_exploit_block"`
	LastAlertTimestamp        interface{} `json:"last_alert_timestamp"`
	LastExploitBlockTimestamp interface{} `json:"last_exploit_block_timestamp"`
	Sysinfo                   struct {
		URL string `json:"url"`
	} `json:"sysinfo"`
	Os struct {
		ProductName   string      `json:"product_name"`
		PatchLevel    interface{} `json:"patch_level"`
		Bitness       string      `json:"bitness"`
		Platform      string      `json:"platform"`
		KernelVersion interface{} `json:"kernel_version"`
	} `json:"os"`
	PrimaryMac string `json:"primary_mac"`
}

type ListOfPolicyResponse struct {
	Details []interface{} `json:"details"`
	Route   string        `json:"route"`
	Data    struct {
		Entries []Policy `json:"entries"`
		Total int `json:"total"`
		Query struct {
			Offset string `json:"offset"`
			Limit  string `json:"limit"`
			Sort   string `json:"sort"`
		} `json:"query"`
		Sort struct {
			Name int `json:"name"`
		} `json:"sort"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"data"`
	Message string `json:"message"`
}

type Proxy struct {

	Host              interface{}   `json:"host"`
	Port              int           `json:"port"`
	Type              string        `json:"type"`
	Enabled           bool          `json:"enabled"`
	Password          interface{}   `json:"password"`
	Username          interface{}   `json:"username"`
	ExcludeHosts      []interface{} `json:"exclude_hosts"`
	FailedRetryDelay  int           `json:"failed_retry_delay"`
	ExcludeLocalHosts bool          `json:"exclude_local_hosts"`
	
}

type Polling struct {
	PollAgentsSecs        int  `json:"poll_agents_secs"`
	ConfigPullEnabled     bool `json:"config_pull_enabled"`
	FastpollAgentsSecs    int  `json:"fastpoll_agents_secs"`
	RequestSysinfoSecs    int  `json:"request_sysinfo_secs"`
	ConfigPollIntervalSec int  `json:"config_poll_interval_sec"`
}

type AgentLogging struct {
	Enabled  bool        `json:"enabled"`
	LogMask  interface{} `json:"log_mask"`
	LogLevel string      `json:"log_level"`
}

type ResourceUse struct {
	Priority                   string `json:"priority"`
	CPULimit                   int    `json:"cpu_limit"`
	MaxDbSize                  int    `json:"max_db_size"`
	StorageMode                string `json:"storage_mode"`
	ConcurrentHostLimit        int    `json:"concurrent_host_limit"`
	ConcurrentHostLimitEnabled bool   `json:"concurrent_host_limit_enabled"`
}

type MalwareScans struct {
	Enable           bool          `json:"enable"`
	ScheduleList     []interface{} `json:"scheduleList"`
	PauseEnabled     bool          `json:"pause_enabled"`
	CancelEnabled    bool          `json:"cancel_enabled"`
	MaxPauseTimes    int           `json:"max_pause_times"`
	MaxPauseDuration int           `json:"max_pause_duration"`
}

type ServerAddress struct {
	Servers []struct {
		Hostname     string `json:"hostname"`
		Provisioning struct {
			Enable        bool `json:"enable"`
			LegacyPrimary bool `json:"legacy_primary"`
		} `json:"provisioning"`
	} `json:"servers"`
}

type RemovalProtection struct {
	Password          interface{} `json:"password"`
	ProtectionEnabled bool        `json:"protection_enabled"`
}

type Policy struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	PolicyTypeID string `json:"policy_type_id"`
	Priority     int    `json:"priority"`
	Enabled      bool   `json:"enabled"`
	Default      bool   `json:"default"`
	Migrated     bool   `json:"migrated"`
	CreatedBy    string `json:"created_by"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Categories   struct {
		Proxy Proxy `json:"proxy"`
		Polling Polling `json:"polling"`
		AgentLogging AgentLogging `json:"agentLogging"`
		ResourceUse ResourceUse `json:"resource_use"`
		MalwareScans MalwareScans `json:"malware_scans"`
		ServerAddress ServerAddress `json:"server_address"`
		MalwareProtection MalwareProtection `json:"malware_protection"`
		RemovalProtection RemovalProtection `json:"removal_protection"`
		ExploitGuardProtection ExploitGuardProtection `json:"exploit_guard_protection"`
		RealTimeIndicatorDetection RealTimeIndicatorDetection `json:"real_time_indicator_detection"`
	} `json:"categories"`
	DisplayCreatedAt string `json:"display_created_at"`
	DisplayUpdatedAt string `json:"display_updated_at"`
}

type MalwareProtection struct {
	
	Enable                        bool        `json:"enable"`
	ExcludedMD5S                  interface{} `json:"excludedMD5s"`
	ExcludedFiles                 interface{} `json:"excludedFiles"`
	UpdateSource                  string      `json:"update_source"`
	ExceptionsPup                 bool        `json:"exceptions_pup"`
	UpdateInterval                int         `json:"update_interval"`
	NetworkOasMode                string      `json:"network_oas_mode"`
	ExceptionsAdware              bool        `json:"exceptions_adware"`
	ExcludedProcesses             interface{} `json:"excludedProcesses"`
	QuarantineEnable              bool        `json:"quarantine_enable"`
	ExceptionsSpyware             bool        `json:"exceptions_spyware"`
	ActionsNotifyUser             bool        `json:"actions_notify_user"`
	NetworkOasEnabled             bool        `json:"network_oas_enabled"`
	ActionsRemoveTrace            bool        `json:"actions_remove_trace"`
	MalOasArchiveScan             bool        `json:"mal_oas_archive_scan"`
	MalOasScanTimeout             int         `json:"mal_oas_scan_timeout"`
	MalOdsArchiveScan             bool        `json:"mal_ods_archive_scan"`
	MalOdsScanTimeout             int         `json:"mal_ods_scan_timeout"`
	MalwareGuardEnable            bool        `json:"malware_guard_enable"`
	QuarantineAgingLimit          int         `json:"quarantine_aging_limit"`
	ActionsCleanInfection         bool        `json:"actions_clean_infection"`
	MalOasAvMaxFileSize           int         `json:"mal_oas_av_max_file_size"`
	MalOasBlockingTimeout         int         `json:"mal_oas_blocking_timeout"`
	MalOasMgMaxFileSize           int         `json:"mal_oas_mg_max_file_size"`
	MalOdsAvMaxFileSize           int         `json:"mal_ods_av_max_file_size"`
	MalOdsMgMaxFileSize           int         `json:"mal_ods_mg_max_file_size"`
	MalOasMaxArchiveDepth         int         `json:"mal_oas_max_archive_depth"`
	MalOdsMaxArchiveDepth         int         `json:"mal_ods_max_archive_depth"`
	ExceptionsHeuristicDetections bool        `json:"exceptions_heuristic_detections"`
	MalwareGuardQuarantineEnable  bool        `json:"malware_guard_quarantine_enable"`
	
}

type ExploitGuardProtection struct {
	ExcludedMD5S       interface{}  `json:"excludedMD5s"`
	ExcludedFiles      interface{}  `json:"excludedFiles"`
	ExcludedPaths      interface{}  `json:"excludedPaths"`
	AlertThreshold     int         `json:"alertThreshold"`
	EnablePageguard    bool        `json:"enable_pageguard"`
	EnableServerOs     bool        `json:"enable_server_os"`
	EnableProduction   bool        `json:"enable_production"`
	EnableTermination  bool        `json:"enable_termination"`
	EnableNotification bool        `json:"enable_notification"`
	EnablePreventKnown bool        `json:"enable_prevent_known"`
}

type RealTimeIndicatorDetection struct {
	ExcludedPaths           interface{} `json:"excludedPaths"`
	IntelPollSec            int         `json:"intel_poll_sec"`
	UDPSendEvents           bool        `json:"udp_send_events"`
	ExcludedProcessNames    interface{} `json:"excludedProcessNames"`
	ActiveCollectionEnabled bool        `json:"active_collection_enabled"`
}

type PolicyUpdate struct {
	Details []interface{} `json:"details"`
	Route   string        `json:"route"`
	Data    []Policy `json:"data"`
	Message string `json:"message"`
}