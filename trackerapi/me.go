package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"
	"strings"

	"github.com/jasondemps1/clirescue/cmdutil"
	"github.com/jasondemps1/clirescue/user"
)

var (
	URL          = "https://www.pivotaltracker.com/services/v5/me"
	FileLocation = homeDir() + "/.tracker"
	currentUser  = user.New()
	Stdout       = os.Stdout
)

func Me() {
	file, err := ioutil.ReadFile(FileLocation)

	if err == nil {
		fmt.Println("Found credentials")

		toStr := string(file[:])

		creds := strings.Split(toStr, ":")

		if len(creds) >= 2 && creds[0] != "" && creds[1] != "" {
			currentUser.Username = creds[0]
			currentUser.Password = creds[1]

			fmt.Printf("Username: %s\nPassword: %s\n", creds[0], creds[1])
		} else {
			fmt.Println("Could not parse credentials file. Please input them again.")
			setCredentials()
		}
	} else {
		setCredentials()
	}

	parse(makeRequest())
	err = writeInfo()

	if err != nil {
		fmt.Printf("Error writing to: %s\n", FileLocation)
	} else {
		fmt.Printf("Successfully Logged In as: %s\n", currentUser.Username)
	}
}

func writeInfo() error {
	return ioutil.WriteFile(FileLocation, []byte(currentUser.Username+":"+currentUser.Password), 0644)
}

func makeRequest() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(currentUser.Username, currentUser.Password)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) error {
	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}
	if meResp.Kind == "error" {
		fmt.Printf("%s\n%s\n", meResp.Error, meResp.PossibleFix)
	}

	currentUser.APIToken = meResp.APIToken
	return err
}

func setCredentials() {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type MeResponse struct {
	APIToken    string `json:"api_token"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Initials    string `json:"initials"`
	Kind        string `json:"kind"`
	Error       string `json:"error"`
	PossibleFix string `json:"possible_fix"`
	Timezone    struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}
