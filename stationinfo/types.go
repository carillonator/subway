package main

type StationInfo struct {
	Id        int
	ComplexId int
	GtfsId    string
	Division  string
	Line      string
	Name      string
	Borough   string
	DayRoutes []string
	Structure string
	Lat       float64
	Long      float64
	Feeds     []uint8
}

type StationInfoSet map[string]*StationInfo

type ComplexInfo struct {
	Id       int
	Name     string
	Stations []*StationInfo
}

type ComplexInfoSet map[int]*ComplexInfo
