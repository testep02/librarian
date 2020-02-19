package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/splitio/go-client/splitio/client"
	"github.com/splitio/go-client/splitio/conf"
	"github.com/splitio/go-toolkit/logging"
)

type eagle struct {
	XMLName xml.Name `xml:"eagle"`
	Drawing drawing  `xml:"drawing"`
}

type drawing struct {
	XMLName  xml.Name `xml:"drawing"`
	Settings settings `xml:"settings"`
	Grid     grid     `xml:"grid"`
	Layers   layers   `xml:"layers"`
	Library  library  `xml:"library"`
}

type settings struct {
	XMLName xml.Name  `xml:"settings"`
	Setting []setting `xml:"setting"`
}

type grid struct {
	XMLName     xml.Name `xml:"grid"`
	Distance    string   `xml:"distance,attr"`
	UnitDist    string   `xml:"unitdist,attr"`
	Unit        string   `xml:"unit,attr"`
	Style       string   `xml:"style,attr"`
	Multiple    string   `xml:"multiple,attr"`
	Display     string   `xml:"display,attr"`
	AltDistance string   `xml:"altdistance,attr"`
	AltUnitDist string   `xml:"altunitdist,attr"`
	AltUnit     string   `xml:"altunit,attr"`
}

type setting struct {
	XMLName           xml.Name `xml:"setting"`
	AlwaysVectorFont  string   `xml:"alwaysvectorfont,attr"`
	KeepOldVectorFont string   `xml:"keepoldvectorfont,attr"`
	VerticalText      string   `xml:"verticaltext,attr"`
}

type layers struct {
	XMLName xml.Name `xml:"layers"`
	Layer   []layer  `xml:"layer"`
}

type layer struct {
	XMLName xml.Name `xml:"layer"`
	Number  string   `xml:"number,attr"`
	Name    string   `xml:"name,attr"`
	Color   string   `xml:"color,attr"`
	Fill    string   `xml:"fill,attr"`
	Visible string   `xml:"visible,attr"`
	Active  string   `xml:"active,attr"`
}

type library struct {
	XMLName    xml.Name   `xml:"library"`
	Devicesets devicesets `xml:"devicesets"`
	Packages   packages   `xml:"packages"`
}

type packages struct {
	XMLName xml.Name     `xml:"packages"`
	Package []libPackage `xml:"package"`
}

type libPackage struct {
	XMLName   xml.Name    `xml:"package"`
	Name      string      `xml:"name,attr"`
	SMD       []smd       `xml:"smd"`
	Text      []text      `xml:"text"`
	Wire      []wire      `xml:"wire"`
	Circle    []circle    `xml:"circle"`
	Polygon   []polygon   `xml:"polygon"`
	Rectangle []rectangle `xml:"rectangle"`
}

type smd struct {
	XMLName xml.Name `xml:"smd"`
	Name    string   `xml:"name,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	DX      string   `xml:"dx,attr"`
	DY      string   `xml:"dy,attr"`
	Layer   string   `xml:"layer,attr"`
}

type symbol struct {
	XMLName     xml.Name    `xml:"symbol"`
	Description string      `xml:"description"`
	Name        string      `xml:"name,attr"`
	Wire        []wire      `xml:"wire"`
	Polygon     []polygon   `xml:"polygon"`
	Text        []text      `xml:"text"`
	Pin         []pin       `xml:"pin"`
	Rectangle   []rectangle `xml:"rectangle"`
	Circle      []circle    `xml:"circle"`
}

type pin struct {
	XMLName  xml.Name `xml:"pin"`
	Name     string   `xml:"name,attr"`
	X        string   `xml:"x,attr"`
	Y        string   `xml:"y,attr"`
	Length   string   `xml:"length,attr"`
	Rotation string   `xml:"rot,attr"`
}

type wire struct {
	XMLName xml.Name `xml:"wire"`
	X1      string   `xml:"x1,attr"`
	Y1      string   `xml:"y1,attr"`
	X2      string   `xml:"x2,attr"`
	Y2      string   `xml:"y2,attr"`
	Width   string   `xml:"width,attr"`
	Layer   string   `xml:"layer,attr"`
}

type text struct {
	XMLName xml.Name `xml:"text"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Size    string   `xml:"size,attr"`
	Layer   string   `xml:"layer,attr"`
}

type circle struct {
	XMLName xml.Name `xml:"circle"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Radius  string   `xml:"radius,attr"`
	Width   string   `xml:"width,attr"`
	Layer   string   `xml:"layer,attr"`
}

type rectangle struct {
	XMLName xml.Name `xml:"rectangle"`
	X1      string   `xml:"x1,attr"`
	Y1      string   `xml:"y1,attr"`
	X2      string   `xml:"x2,attr"`
	Y2      string   `xml:"y2,attr"`
	Layer   string   `xml:"layer,attr"`
}

type polygon struct {
	XMLName xml.Name `xml:"polygon"`
	Width   string   `xml:"width,attr"`
	Layer   string   `xml:"layer,attr"`
	Vertex  []vertex `xml:"vertex"`
}

type vertex struct {
	XMLName xml.Name `xml:"vertex"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
}

type devicesets struct {
	XMLName   xml.Name    `xml:"devicesets"`
	Deviceset []deviceset `xml:"deviceset"`
}

type deviceset struct {
	XMLName     xml.Name `xml:"deviceset"`
	Name        string   `xml:"name,attr"`
	Prefix      string   `xml:"prefix,attr"`
	UserValue   string   `xml:"uservalue,attr"`
	Description string   `xml:"description"`
	Gates       gates    `xml:"gates"`
	Devices     devices  `xml:"devices"`
}

type gates struct {
	XMLName xml.Name `xml:"gates"`
	Gate    []gate   `xml:"gate"`
}

type gate struct {
	XMLName xml.Name `xml:"gate"`
	Name    string   `xml:"name,attr"`
	Symbol  string   `xml:"symbol,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
}

type devices struct {
	XMLName xml.Name `xml:"devices"`
	Device  []device `xml:"device"`
}

type device struct {
	XMLName      xml.Name     `xml:"device"`
	Connects     connects     `xml:"connects"`
	Name         string       `xml:"name,attr"`
	Package      string       `xml:"package,attr"`
	Technologies technologies `xml:"technologies"`
}

type connects struct {
	XMLName xml.Name  `xml:"connects"`
	Connect []connect `xml:"connect"`
}

type connect struct {
	XMLName xml.Name `xml:"connect"`
	Gate    string   `xml:"gate,attr"`
	Pin     string   `xml:"pin,attr"`
	Pad     string   `xml:"pad,attr"`
}

type technologies struct {
	XMLName    xml.Name     `xml:"technologies"`
	Technology []technology `xml:"technology"`
}

type technology struct {
	XMLName       xml.Name        `xml:"technology"`
	Name          string          `xml:"name,attr"`
	TechAttribute []techAttribute `xml:"attribute"`
}

type techAttribute struct {
	XMLName  xml.Name `xml:"attribute"`
	Name     string   `xml:"name,attr"`
	Value    string   `xml:"value,attr"`
	Constant string   `xml:"constant,attr"`
}

func main() {

	/* This section begins setup for Split.io Feature Flagging */
	cfg := conf.Default()                          // This creates a new splitio configuration that will be passed to the client
	cfg.LoggerConfig.LogLevel = logging.LevelDebug // Set the logging level that will be output to the console

	// This line of code creates a new splitio client and passes in the config created above
	factory, err := client.NewSplitFactory("gjfmd11c290hph9j73bbp5jjpefn6gei5rmf", cfg)

	// Check for any errors that occurred while creating the client
	if err != nil {
		fmt.Printf("SDK Init Error: %s\n", err)
		return
	}

	// Here we grab the actual client from the factory that was created
	client := factory.Client()
	// It can take some time for the client to be ready, so we need to
	// block the application until we are ready for the client to make
	// requests. We can also do this in a channel to prevent blocking
	// the entire app. That would be the preferred way to do it
	err = client.BlockUntilReady(10)

	// Check for any errors in client creation from the factory
	if err != nil {
		fmt.Println("Error getting client:", err)
		return
	}

	/* 	A treatment is a flag
	 	This call will bring back a string value
		Possible values:
			If user is in the QA group: "on"
			If user is not in the QA group: "inBetween"
			If no user is passed in: "off"

		The "inBetween" value shows how the flag can have
		a non-standard on/off value and how that value
		can be tied to attributes or users that are
		passed in.

	*/

	treatment := client.Treatment("", "print_devices", nil)
	// To test this treatment, you can use these values:
	//
	// To evaluate to "on":
	// Set first argument to any of these values:
	//		- Lee Roy
	// 		- Billy Bob
	//		- Someone Else
	// 		- Test User
	//
	// To evaluate to "inBetween":
	// Set first argument to any string other than those
	// listed above.

	eagleFilePath, err := filepath.Abs("common.lbr")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(eagleFilePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	eagleDrawing, err := readEagle(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	devices := eagleDrawing.Library.Devicesets.Deviceset

	/******************************************************************
	*	This section evaluates the feature flag and executes code
	*	based on the value that we get back
	******************************************************************/

	// If the user we pass in is in the QA group:
	if treatment == "on" {
		for d := range devices {
			fmt.Println("Device Name:", devices[d].Name)
		}
		// There are no situations currently where the code
		// will reach here
	} else if treatment == "off" {
		fmt.Println("Device listing is turned off.")
		fmt.Println()
		fmt.Println()
		//If a user is passed in, but they are not in the QA group
	} else if treatment == "inBetween" {
		fmt.Println("Treatment is", treatment)
		fmt.Println()
		fmt.Println()
		// If an empty string is passed in instead of a name
		// this is the code that will be executed
	} else {
		fmt.Println("Device listing is neither on nor off....wtf??")
		fmt.Println()
		fmt.Println()
	}

	/******************************************************************
	*	End feature flagging PoC section
	******************************************************************/

	fmt.Printf("The Layer Name setting is: %s", eagleDrawing.Layers.Layer[0].Name)
	fmt.Println()
	fmt.Printf("The Layer Number setting is: %s", eagleDrawing.Layers.Layer[0].Number)
	fmt.Println()
	fmt.Printf("The Layer Color is: %s", eagleDrawing.Layers.Layer[0].Color)
	fmt.Println()
	fmt.Printf("The Technology Name setting is: %s", eagleDrawing.Library.Devicesets.Deviceset[0].Devices.Device[0].Technologies.Technology[0].Name)
	fmt.Println()
	//fmt.Printf("The Technology Package setting is: %s", eagleDrawing.Library.Devicesets.Deviceset[0].Devices.Device[0].Technologies.Technology[0].TechAttribute[0].Name)
	fmt.Println()
	fmt.Printf("The Layer Color is: %s", eagleDrawing.Layers.Layer[0].Color)
	fmt.Println()
	fmt.Printf("The SMD Name is: %s", eagleDrawing.Library.Packages.Package[0].SMD[0].Name)
	fmt.Println()
	fmt.Printf("The Gate Symbol Name is: %s", eagleDrawing.Library.Devicesets.Deviceset[0].Gates.Gate[0].Symbol)
	fmt.Println()
}

func readEagle(reader io.Reader) (drawing, error) {
	var xmlEagle eagle

	if err := xml.NewDecoder(reader).Decode(&xmlEagle); err != nil {
		return xmlEagle.Drawing, err
	}

	return xmlEagle.Drawing, nil
}
