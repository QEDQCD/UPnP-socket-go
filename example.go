package main
import (
	"fmt"
	"os"
        "net"
        "time"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
)

func main() {
	clients, errors, err := internetgateway1.NewWANPPPConnection1Clients()
	extIPClients := make([]GetExternalIPAddresser, len(clients))
	for i, client := range clients {
		extIPClients[i] = client
	}
	 DisplayExternalIPResults(extIPClients, errors, err)
	 Example_WANIPConnection_GetExternalIPAddress()
}

func Example_WANIPConnection_GetExternalIPAddress() {
	clients, errors, err := internetgateway1.NewWANIPConnection1Clients()
	extIPClients := make([]GetExternalIPAddresser, len(clients))
	for i, client := range clients {
		extIPClients[i] = client
	}
	DisplayExternalIPResults(extIPClients, errors, err)
	// Output:
}

type GetExternalIPAddresser interface {
	GetExternalIPAddress() (NewExternalIPAddress string, err error)
	GetServiceClient() *goupnp.ServiceClient
                GetStatusInfo()(string, string, uint32,error)
                GetIdleDisconnectTime()(uint32, error)
                AddPortMapping (string, uint16, string, uint16, string, bool, string, uint32)error
}
func DisplayExternalIPResults(clients []GetExternalIPAddresser, errors []error, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error discovering service with UPnP: ", err)
		return
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Error discovering %d services:\n", len(errors))
		for _, err := range errors {
			fmt.Println("  ", err)
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully discovered %d services:\n", len(clients))
	for _, client := range clients {
		device := &client.GetServiceClient().RootDevice.Device

		fmt.Fprintln(os.Stderr, "  Device:", device.FriendlyName)
		if addr, err := client.GetExternalIPAddress(); err != nil {
			fmt.Fprintf(os.Stderr, "    Failed to get external IP address: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "    External IP address: %v\n", addr)
		}
                          srv := client.GetServiceClient().Service
		fmt.Println(device.FriendlyName, " :: ", srv.String())
		scpd, err := srv.RequestSCPD()
                             if err != nil {
			fmt.Printf("  Error requesting service SCPD: %v\n", err)
		} else {
			fmt.Println("  Available actions:")
			for _, action := range scpd.Actions {
				fmt.Printf("  * %s\n", action.Name)
				for _, arg := range action.Arguments {
					var varDesc string
					if stateVar := scpd.GetStateVariable(arg.RelatedStateVariable); stateVar != nil {
						varDesc = fmt.Sprintf(" (%s)", stateVar.DataType.Name)
					}
					fmt.Printf("    * [%s] %s%s\n", arg.Direction, arg.Name, varDesc)
				}
			}
                            }
                         if scpd == nil || scpd.GetAction("GetExternalIPAddress") != nil {
			ip, err :=  client.GetExternalIPAddress()
			fmt.Println("GetExternalIPAddress: ", ip, err)
		}

                            if scpd == nil || scpd.GetAction("GetStatusInfo") != nil {
			status, lastErr, uptime, err := client.GetStatusInfo()
			fmt.Println("GetStatusInfo: ", status, lastErr, uptime, err)
		}

		if scpd == nil || scpd.GetAction("GetIdleDisconnectTime") != nil {
			idleTime, err := client.GetIdleDisconnectTime()
			fmt.Println("GetIdleDisconnectTime: ", idleTime, err)
		}

		if scpd == nil || scpd.GetAction("AddPortMapping") != nil  {
			err :=  client.AddPortMapping("", 5000, "TCP", 5001, "192.168.1.20", true, "Test port mapping", 0)
			fmt.Println("AddPortMapping: ", err)
		}
                      testsocket()
	}
}
func testsocket() {
server := "192.168.1.20:5001"
    netListen, err := net.Listen("tcp", server)
    if err != nil{
        Log("connect error: ", err)
        os.Exit(1)
    }
    Log("Waiting for Client ...")
    for{ 
        conn, err := netListen.Accept()   
        if err != nil{
            Log(conn.RemoteAddr().String(), "Fatal error: ", err)
            continue
        }
 
       
        conn.SetReadDeadline(time.Now().Add(time.Duration(10)*time.Second))
 
        Log(conn.RemoteAddr().String(), "connect success!")
        go handleConnection(conn)
 
    }
}
func handleConnection(conn net.Conn) {
    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if err != nil {     
            return
        }
 
        Data := buffer[:n]
        message := make(chan byte)
      
        go HeartBeating(conn, message, 4)
        
        go GravelChannel(Data, message)
 
        Log(time.Now().Format("2006-01-02 15:04:05.0000000"), conn.RemoteAddr().String(), string(buffer[:n]))
    }
 
    defer conn.Close()
}
func GravelChannel(bytes []byte, mess chan byte) {
    for _, v := range bytes{
        mess <- v
    }
    close(mess)
}
func HeartBeating(conn net.Conn, bytes chan byte, timeout int) {
    select {
    case fk := <- bytes:
        Log(conn.RemoteAddr().String(), "heartbeat", string(fk), "times")
        conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
        break
 
        case <- time.After(5 * time.Second):
            Log("conn dead now")
            conn.Close()
    }
}
func Log(v ...interface{}) {
    fmt.Println(v...)
    return
}
